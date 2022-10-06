package main

import (
	"fmt"
	"math"

	"github.com/tuneinsight/lattigo/v3/ckks"
	"github.com/tuneinsight/lattigo/v3/dckks"
	"github.com/tuneinsight/lattigo/v3/drlwe"
	"github.com/tuneinsight/lattigo/v3/rlwe"
	"github.com/tuneinsight/lattigo/v3/utils"
)

type party struct {
	sk *rlwe.SecretKey
	//rlkEphemSk *rlwe.SecretKey

	ckgShare *drlwe.CKGShare
	//rkgShareOne *drlwe.RKGShare
	//rkgShareTwo *drlwe.RKGShare
	pcksShare [][]*drlwe.PCKSShare // PKCS protocol -> public key switching protocol

	input  [][]float64 //uint64 // fingerprint
	NumRow int         // NumCol real
	NumCol int         // NumRow real
}

func getParameters(prnus [][][]PixelGray, N int) [][]float64 {

	// EN UN CIPHERTEXT:
	// Con CKKS -> 2^LogSlots (aquí 2048 valores)
	// Con BFV -> 2^logN

	// ENCRYPTION PARAMETERS

	paramsDef := ckks.PN12QP109 // default --> LogSlots = 11 => 2^11 = 2048 valores en el ciphertext
	//paramsDef.T = 1
	params, err := ckks.NewParametersFromLiteral(paramsDef)
	if err != nil {
		panic(err)
	}

	crs, err := utils.NewKeyedPRNG([]byte{'f', 'e', 'l', 'd', 's', 'p', 'a', 'r'}) //'t', 'r', 'u', 'm', 'p', 'e', 't'
	if err != nil {
		panic(err)
	}

	encoder := ckks.NewEncoder(params)

	// Target private and public keys
	tsk, tpk := ckks.NewKeyGenerator(params).GenKeyPair()

	// Create each party and allocate the memory for all the shares that the protocols will need
	P := genparties(params, N)

	// Inputs & expected result
	//expRes := genInputs(params, P, OrNumRow, OrNumCol, 0xffffffffffffffff)

	// Assign inputs (each prnu to each user/party)
	getInputs(P, prnus)

	// 1) Collective public key generation
	pk := ckgphase(params, crs, P)

	// encInputs (get encrypted prnus)
	encInputs := encPhase(params, P, pk, encoder)

	encRes := evalPhase(params, encInputs)

	encOut := pcksPhase(params, tpk, encRes, P)
	fmt.Printf("Size of result\t: NumRow: %d ciphertexts, NumCol: %d ciphertexts\n", len(encOut), len(encOut[0]))

	// Decrypt the result with the target secret key
	fmt.Println("> Decrypt Phase")
	decryptor := ckks.NewDecryptor(params, tsk)

	ptres := make([][]*ckks.Plaintext, len(encOut))
	for i := range encOut {
		ptres[i] = make([]*ckks.Plaintext, len(encOut[i]))
		for j := range encOut[i] {
			ptres[i][j] = ckks.NewPlaintext(params, 1, params.DefaultScale())
		}
	}

	for i := range encOut {
		for j := range encOut[i] {
			decryptor.Decrypt(encOut[i][j], ptres[i][j])
		}
	}
	fmt.Println("done")

	fmt.Println("> Result:")
	// Check the result
	res := make([][]float64, P[0].NumRow)
	for i := range ptres {
		res[i] = make([]float64, P[0].NumCol)
	}
	//fmt.Printf("size res: Row %d x Col %d\n", len(ptres), len(ptres[0]))
	//fmt.Printf("size res: Row %d x Col %d\n", len(res), len(res[0]))

	for i := range ptres {
		for j := range ptres[i] {
			partialRes := encoder.DecodeCoeffs(ptres[i][j])
			for k := range partialRes {
				res[i][(j*len(partialRes) + k)] = partialRes[k]
			}
		}
	}

	fmt.Println("finish")

	return res
}

// Generates the invidividual secret key and "input images of size Num_Row x Num_Col" for each Forensic Party P[i]
func genparties(params ckks.Parameters, N int) []*party { //genSKparties

	// Create each party, and allocate the memory for all the shares that the protocols will need
	P := make([]*party, N)
	for i := range P {
		pi := &party{}
		pi.sk = ckks.NewKeyGenerator(params).GenSecretKey()
		P[i] = pi
		//P[i].sk = ckks.NewKeyGenerator(params).GenSecretKey()
	}

	return P
}

func getInputs(p []*party, prnus [][][]PixelGray) {

	in := make([][]float64, len(prnus[0]))

	for t := 0; t < len(in); t++ {
		in[t] = make([]float64, len(prnus[0][0]))
	}

	// prnus ó P => [user][rows][columns]
	for i := 0; i < len(p); i++ { // len(P) = len(prnus)

		p[i].input = in

		for j := 0; j < len(prnus[i]); j++ {
			for k := 0; k < len(prnus[i][j]); k++ {

				p[i].input[j][k] = prnus[i][j][k].pix
			}
		}

		p[i].NumRow = len(p[i].input)
		fmt.Println(p[i].NumRow)
		p[i].NumCol = len(p[i].input[0])
		fmt.Println(p[i].NumCol)

	}

}

func ckgphase(params ckks.Parameters, crs utils.PRNG, P []*party) *rlwe.PublicKey {

	ckg := dckks.NewCKGProtocol(params) // Public key generation
	ckgCombined := ckg.AllocateShare()

	for _, pi := range P {
		pi.ckgShare = ckg.AllocateShare()
	}

	crp := ckg.SampleCRP(crs)

	for _, pi := range P {
		ckg.GenShare(pi.sk, crp, pi.ckgShare)
	}

	pk := ckks.NewPublicKey(params)

	for _, pi := range P {
		ckg.AggregateShare(pi.ckgShare, ckgCombined, ckgCombined)
	}
	ckg.GenPublicKey(ckgCombined, crp, pk)

	return pk
}

func encPhase(params ckks.Parameters, P []*party, pk *rlwe.PublicKey, encoder ckks.Encoder) (encInputs [][][]*ckks.Ciphertext) {

	NumRowEncIn := P[0].NumRow                                                    // 5
	NumColEncIn := int(math.Ceil(float64(P[0].NumCol) / float64(params.Slots()))) // numCol/2048
	// ceil => redondeo hacia arriba (igual es necesario rellenar con ceros)

	// SIZE OF THE CIPHERTEXT: 2048 valores

	// encInputs[i][j][k], i through Parties, j through Rows, k through Columns
	// empty ciphertext
	encInputs = make([][][]*ckks.Ciphertext, len(P))
	for i := range encInputs {
		encInputs[i] = make([][]*ckks.Ciphertext, NumRowEncIn)
		for j := range encInputs[i] {
			encInputs[i][j] = make([]*ckks.Ciphertext, NumColEncIn)
		}
	}

	//encOut[i] = ckks.NewCiphertext(params, encRes[0].Degree(), encRes[0].Level(), encRes[0].Scale)

	// Initializes "input" ciphertexts
	// put info in empty ciphertexts
	for i := range encInputs {
		fmt.Printf("i = %d\n", i)
		for j := range encInputs[i] {
			fmt.Printf("j = %d\n", j)
			for k := range encInputs[i][j] {
				fmt.Printf("k = %d\n", k)
				encInputs[i][j][k] = ckks.NewCiphertext(params, 1 /*int(params.N())*/, 1, params.DefaultScale())
			}
		}
	}

	// ENCRYPT PHASE
	// Each party encrypts its bidimensional array of input vectors into a bidimensional array of input ciphertexts
	encryptor := ckks.NewEncryptor(params, pk)

	pt := ckks.NewPlaintext(params, 1, params.DefaultScale())

	// create cyphertexts
	for i, pi := range P {
		for j := 0; j < NumRowEncIn; j++ {
			for k := 0; k < NumColEncIn; k++ {

				//rellenar con ceros el ciphertext (si es más grande que los elementos que quedan por cifrar)
				if (k+1)*params.Slots() > len(pi.input[j]) {
					fmt.Println(k)
					fmt.Println((k + 1) * params.Slots())
					fmt.Println((k+1)*params.Slots() - len(pi.input[j]))
					add := make([]float64, (k+1)*params.Slots()-len(pi.input[j])) // slice de ceros
					fmt.Println(add)

					pi.input[j] = append(pi.input[j], add...)
					fmt.Println("hola")
				}

				fmt.Printf("SIZE EACH ROW: %d\n", len(pi.input[j][(k*params.Slots()):((k+1)*params.Slots()-1)]))
				fmt.Printf("valores son %d y %d\n", k*params.Slots(), (k+1)*params.Slots())
				fmt.Printf("size total row %d\n", len(pi.input[j]))

				encoder.Encode(pi.input[j][(k*params.Slots()):((k+1)*params.Slots())], pt, params.LogSlots()) // es uno menos en la segunda parte porque go indexa [0:n] como los valores 0, 1, ..., n - 1
				encryptor.Encrypt(pt, encInputs[i][j][k])
			}
		}
	}

	return
}

// !! modify to include the operation I am interested
// matching case requires. (1) rlk, (2) extra dimension in encRes
func evalPhase(params ckks.Parameters, encInputs [][][]*ckks.Ciphertext) (encRes [][]*ckks.Ciphertext) {

	// Rows, Cols for the "matrices of ciphertexts"
	NumRowEncIn := len(encInputs[0])
	NumColEncIn := len(encInputs[0][0])

	encRes = make([][]*ckks.Ciphertext, NumRowEncIn)
	for i := 0; i < len(encRes); i++ {
		encRes[i] = make([]*ckks.Ciphertext, NumColEncIn)
	}

	for i := 0; i < len(encRes); i++ {
		for j := 0; j < len(encRes[0]); j++ {
			encRes[i][j] = ckks.NewCiphertext(params, 1, 1, params.DefaultScale())
		}
	}

	evaluator := ckks.NewEvaluator(params, rlwe.EvaluationKey{Rlk: nil, Rtks: nil})

	for i := 0; i < len(encInputs); i++ {
		for j := 0; j < len(encInputs[0]); j++ { // NumRowEncIn
			for k := 0; k < len(encInputs[0][0]); k++ { // NumColEncIn
				evaluator.Add(encRes[j][k], encInputs[i][j][k], encRes[j][k])
			}
		}
	}

	/*
		// Layers definition to store intermediate and final results
		encLayers := make([][][][]*ckks.Ciphertext, 0) // array of an array with matrices of ciphertexts
		encLayers = append(encLayers, encInputs)

		for nLayer := len(encInputs) / 2; nLayer > 0; nLayer = nLayer >> 1 {
			encLayer := make([][][]*ckks.Ciphertext, nLayer) // one layer with several matrices of ciphertexts
			for i := range encLayer {                        // Running through ciphertexts in one level
				encLayer[i] = make([][]*ckks.Ciphertext, NumRowEncIn)
				for j := range encLayer[i] { // Running through rows of each matrix of ciphertexts
					encLayer[i][j] = make([]*ckks.Ciphertext, NumColEncIn)
					for k := range encLayer[i][j] { // Running through columns of each row of a matrix of ciphertexts
						encLayer[i][j][k] = ckks.NewCiphertext(params, 1, 1, params.DefaultScale()) // Change in the SECOND EXAMPLE to "bfv.NewCiphertext(params, 2)" to store the result of one multiplication
					}
				}
			}
			encLayers = append(encLayers, encLayer)
		}
		encRes = encLayers[len(encLayers)-1][0]

		//
		evaluator := ckks.NewEvaluator(params, rlwe.EvaluationKey{Rlk: nil, Rtks: nil}) // REMOVE Rlk - SECOND EXAMPLE -> if using evaluator.innersum, we have to generate the power-of-two rotations
		// generar las rotation matrices con los automorfismos de GaloisElementsForRowInnerSum()

		for i := 0; i < len(encRes); i++ {
			for j := 0; i < len(encRes[0]); j++ {
				evaluator.Add(encRes[i][j], encRes[i][j], encRes[i][j]) // SECOND EXAMPLE -> chambiar por mult + relinearize + innersum -> opción de mejora con extractLWEsample
			}
		}*/

	return
}

// cambio de la global secret key a la "target secret key"
func pcksPhase(params ckks.Parameters, tpk *rlwe.PublicKey, encRes [][]*ckks.Ciphertext, P []*party) (encOut [][]*ckks.Ciphertext) {

	// Collective key switching from the collective secret key to
	// the target public key

	//CHECK -> encOut and encRes are matrices of ciphertexts now
	pcks := dckks.NewPCKSProtocol(params, 3.19)

	for _, pi := range P {
		pi.pcksShare = make([][]*drlwe.PCKSShare, len(encRes))
		for i := range encRes {
			pi.pcksShare[i] = make([]*drlwe.PCKSShare, len(encRes[i]))
			for j := range encRes[0] {
				pi.pcksShare[i][j] = pcks.AllocateShare(encRes[0][0].Level())
			}
		}
	}

	for _, pi := range P {
		for i := range encRes {
			for j := range encRes[0] {
				pcks.GenShare(pi.sk, tpk, encRes[i][j].Value[1], pi.pcksShare[i][j])
			}
		}
	}

	//var pcksCombined [][]*drlwe.PCKSShare
	pcksCombined := make([][]*drlwe.PCKSShare, len(encRes))
	encOut = make([][]*ckks.Ciphertext, len(encRes))
	for i := range encRes {
		pcksCombined[i] = make([]*drlwe.PCKSShare, len(encRes[i]))
		encOut[i] = make([]*ckks.Ciphertext, len(encRes[i]))
		for j := range encRes[0] {
			pcksCombined[i][j] = pcks.AllocateShare(encRes[0][0].Level())
			encOut[i][j] = ckks.NewCiphertext(params, 1, 1, params.DefaultScale())
		}
	}

	for _, pi := range P {
		for i := range encRes {
			for j := range encRes[0] {
				pcks.AggregateShare(pi.pcksShare[i][j], pcksCombined[i][j], pcksCombined[i][j])
			}
		}
	}

	for i := range encRes {
		for j := range encRes[0] {
			pcks.KeySwitch(encRes[i][j], pcksCombined[i][j], encOut[i][j])
		}
	}

	return

}
