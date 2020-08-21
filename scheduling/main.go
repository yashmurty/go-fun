package main

import "fmt"

type PatientRawData struct {
	patientID         int
	posTime           int
	irrTime           int
	totalTime         int
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

	simulateAppointments(patientRawData)

	printDataInTable(patientRawData)

}

func simulateAppointments(patientRawData []PatientRawData) {

	availableTimePosRoom := 0
	availableTimeIrrRoom := 0

	for i := 0; i < len(patientRawData); i++ {
		patientDataEach := &patientRawData[i]

		if availableTimeIrrRoom-patientDataEach.posTime >= availableTimePosRoom {
			availableTimePosRoom = availableTimeIrrRoom - patientDataEach.posTime
		}

		patientDataEach.posStartTime = availableTimePosRoom
		patientDataEach.posEndTime = patientDataEach.posStartTime + patientDataEach.posTime

		// NOTE: Calculate final Pos Room visit time for previous patient.
		if i > 0 {
			patientPrevious := &patientRawData[i-1]
			patientPrevious.posFinalStartTime = patientDataEach.posEndTime
			patientPrevious.posFinalEndTime = patientPrevious.posFinalStartTime + FinalVisitPosRoomTime
		}

		if availableTimeIrrRoom <= patientDataEach.posEndTime {
			availableTimeIrrRoom = patientDataEach.posEndTime
		}

		patientDataEach.irrStartTime = availableTimeIrrRoom
		patientDataEach.irrEndTime = patientDataEach.irrStartTime + patientDataEach.irrTime

		// Note: Compensating here for final visit as well.
		availableTimePosRoom = availableTimePosRoom + patientDataEach.posTime + FinalVisitPosRoomTime
		availableTimeIrrRoom = availableTimeIrrRoom + patientDataEach.irrTime

		fmt.Printf("patientDataEach : %+v\n", patientDataEach)

		// break
	}

	fmt.Println("availableTimePosRoom : ", availableTimePosRoom)
	fmt.Println("availableTimeIrrRoom : ", availableTimeIrrRoom)

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

func printDataInTable(patientRawData []PatientRawData) {
	fmt.Println("patientRawData : ", patientRawData)

	fmt.Printf("---- : --- | Pos Room  | Irr Room  | \n")

	for i := 0; i < 170; i = i + 5 {
		fmt.Printf("Time : %3d | ", i)
		patientExistsInPosRoom := false
		patientExistsInIrrRoom := false

		for _, patientRawDataEach := range patientRawData {
			if (patientRawDataEach.posStartTime <= i && patientRawDataEach.posEndTime > i) || (patientRawDataEach.posFinalStartTime <= i && patientRawDataEach.posFinalEndTime > i) {
				fmt.Printf("Patient %1d | ", patientRawDataEach.patientID)
				patientExistsInPosRoom = true

				break
			}
		}
		if !patientExistsInPosRoom {
			fmt.Printf("--------- | ")
		}

		for _, patientRawDataEach := range patientRawData {
			if patientRawDataEach.irrStartTime <= i && patientRawDataEach.irrEndTime > i {
				fmt.Printf("Patient %1d |", patientRawDataEach.patientID)
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
