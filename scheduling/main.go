package main

import "fmt"

// PatientRawData represents the raw data received via CSV files.
type PatientRawData struct {
	patientID int
	posTime   int
	irrTime   int
	totalTime int
}

// PatientAppointmentData stores the appointment schedule timing.
type PatientAppointmentData struct {
	patientID         int
	posStartTime      int
	posEndTime        int
	irrStartTime      int
	irrEndTime        int
	posFinalStartTime int
	posFinalEndTime   int
}

// FinalVisitPosRoomTime simulates the second time the patient visits the Pos room after finishing with IRR room.
const FinalVisitPosRoomTime = 5

func main() {
	var patientRawData = getPatientRawData()

	fmt.Println("patientRawData : ", patientRawData)

	_, availableTimePosRoom, availableTimeIrrRoom := getOptimalAppointment(patientRawData)

	fmt.Println("Optimum availableTimePosRoom : ", availableTimePosRoom)
	fmt.Println("Optimum availableTimeIrrRoom : ", availableTimeIrrRoom)

}

func getOptimalAppointment(patientRawData []PatientRawData) ([]PatientAppointmentData, int, int) {

	optimalPatientAppointmentData := []PatientAppointmentData{}
	optimalTimePosRoom := int(^uint(0) >> 1) // Max possible int
	optimalTimeIrrRoom := int(^uint(0) >> 1) // Max possible int

	var permutations = getElementPermutations(len(patientRawData))

	fmt.Println("permutations : ", permutations)
	fmt.Println("permutationsCount : ", len(permutations))

	for i := 0; i < len(permutations); i++ {
		fmt.Printf("\n--------- \n")
		fmt.Println("permutation : ", permutations[i])
		permutatedPatientRawData := getPermutatedPatientRawData(permutations[i], patientRawData)
		fmt.Println("permutatedPatientRawData : ", permutatedPatientRawData)

		patientAppointmentData, availableTimePosRoom, availableTimeIrrRoom := simulateAppointments(permutatedPatientRawData)
		if availableTimePosRoom < optimalTimePosRoom {
			optimalPatientAppointmentData = patientAppointmentData
			optimalTimePosRoom = availableTimePosRoom
			optimalTimeIrrRoom = availableTimeIrrRoom
		}
	}

	return optimalPatientAppointmentData, optimalTimePosRoom, optimalTimeIrrRoom
}

func simulateAppointments(patientRawData []PatientRawData) ([]PatientAppointmentData, int, int) {

	availableTimePosRoom := 0
	availableTimeIrrRoom := 0

	patientAppointmentData := make([]PatientAppointmentData, len(patientRawData))

	for i := 0; i < len(patientRawData); i++ {
		patientRawDataEach := patientRawData[i]
		patientAppointmentDataEach := &patientAppointmentData[i]
		patientAppointmentDataEach.patientID = patientRawDataEach.patientID

		if availableTimeIrrRoom-patientRawDataEach.posTime >= availableTimePosRoom {
			availableTimePosRoom = availableTimeIrrRoom - patientRawDataEach.posTime
		}

		patientAppointmentDataEach.posStartTime = availableTimePosRoom
		patientAppointmentDataEach.posEndTime = patientAppointmentDataEach.posStartTime + patientRawDataEach.posTime

		// NOTE: Calculate final Pos Room visit time for previous patient.
		if i > 0 {
			patientAppointmentDataPrevious := &patientAppointmentData[i-1]

			patientAppointmentDataPrevious.posFinalStartTime = patientAppointmentDataEach.posEndTime
			patientAppointmentDataPrevious.posFinalEndTime = patientAppointmentDataPrevious.posFinalStartTime + FinalVisitPosRoomTime

			// Compensating here for final visit to the Pos room as well.
			availableTimePosRoom = availableTimePosRoom + FinalVisitPosRoomTime
		}

		if availableTimeIrrRoom <= patientAppointmentDataEach.posEndTime {
			availableTimeIrrRoom = patientAppointmentDataEach.posEndTime
		}

		patientAppointmentDataEach.irrStartTime = availableTimeIrrRoom
		patientAppointmentDataEach.irrEndTime = patientAppointmentDataEach.irrStartTime + patientRawDataEach.irrTime

		availableTimePosRoom = availableTimePosRoom + patientRawDataEach.posTime
		availableTimeIrrRoom = availableTimeIrrRoom + patientRawDataEach.irrTime

		// Calculate final Pos Room visit for the final patient.
		if i == len(patientRawData)-1 {
			if availableTimeIrrRoom >= availableTimePosRoom {
				availableTimePosRoom = availableTimeIrrRoom
			}

			patientAppointmentDataEach.posFinalStartTime = availableTimePosRoom
			patientAppointmentDataEach.posFinalEndTime = patientAppointmentDataEach.posFinalStartTime + FinalVisitPosRoomTime

			// Compensating here for final visit to the Pos room as well.
			availableTimePosRoom = availableTimePosRoom + FinalVisitPosRoomTime
		}

		// fmt.Printf("patientAppointmentDataEach : %+v\n", patientAppointmentDataEach)

	}

	fmt.Println("availableTimePosRoom : ", availableTimePosRoom)
	fmt.Println("availableTimeIrrRoom : ", availableTimeIrrRoom)

	printDataInTable(patientAppointmentData)

	return patientAppointmentData, availableTimePosRoom, availableTimeIrrRoom
}

func getPatientRawData() []PatientRawData {
	return []PatientRawData{
		{
			patientID: 1,
			posTime:   5,
			irrTime:   20,
			totalTime: 25,
		},
		{
			patientID: 2,
			posTime:   10,
			irrTime:   15,
			totalTime: 25,
		},
		{
			patientID: 3,
			posTime:   5,
			irrTime:   15,
			totalTime: 20,
		},
		{
			patientID: 4,
			posTime:   10,
			irrTime:   25,
			totalTime: 35,
		},
		{
			patientID: 5,
			posTime:   10,
			irrTime:   15,
			totalTime: 25,
		},
		{
			patientID: 6,
			posTime:   25,
			irrTime:   10,
			totalTime: 35,
		},
		{
			patientID: 7,
			posTime:   10,
			irrTime:   10,
			totalTime: 20,
		},
		{
			patientID: 8,
			posTime:   5,
			irrTime:   10,
			totalTime: 15,
		},
		{
			patientID: 9,
			posTime:   10,
			irrTime:   15,
			totalTime: 25,
		},
	}
}

func getPermutatedPatientRawData(permutationOrder []int, patientRawData []PatientRawData) []PatientRawData {
	permutatedPatientRawData := make([]PatientRawData, len(patientRawData))
	for i := 0; i < len(permutationOrder); i++ {
		permutatedPatientRawData[i] = patientRawData[permutationOrder[i]]
	}

	return permutatedPatientRawData
}

func getElementPermutations(count int) [][]int {
	// WIP: Using manual permutations for now. Will update this later to be automatic.
	permutation1 := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	permutation2 := []int{1, 0, 2, 3, 4, 5, 6, 7, 8}

	permutations := [][]int{permutation1, permutation2}
	return permutations
}

func printDataInTable(patientAppointmentData []PatientAppointmentData) {
	fmt.Println("patientAppointmentData : ", patientAppointmentData)

	fmt.Printf("---- : --- | Pos Room  | Irr Room  | \n")

	for i := 0; i < 185; i = i + 5 {
		fmt.Printf("Time : %3d | ", i)
		patientExistsInPosRoom := false
		patientExistsInIrrRoom := false

		for _, patientAppointmentDataEach := range patientAppointmentData {
			if (patientAppointmentDataEach.posStartTime <= i && patientAppointmentDataEach.posEndTime > i) || (patientAppointmentDataEach.posFinalStartTime <= i && patientAppointmentDataEach.posFinalEndTime > i) {
				fmt.Printf("Patient %1d | ", patientAppointmentDataEach.patientID)
				patientExistsInPosRoom = true

				break
			}
		}
		if !patientExistsInPosRoom {
			fmt.Printf("--------- | ")
		}

		for _, patientAppointmentDataEach := range patientAppointmentData {
			if patientAppointmentDataEach.irrStartTime <= i && patientAppointmentDataEach.irrEndTime > i {
				fmt.Printf("Patient %1d |", patientAppointmentDataEach.patientID)
				patientExistsInIrrRoom = true

				break
			}
		}
		if !patientExistsInIrrRoom {
			fmt.Printf("--------- |")
		}

		fmt.Printf(" \n")
	}

}
