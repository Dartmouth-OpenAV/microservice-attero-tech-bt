package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Dartmouth-OpenAV/microservice-framework/framework"
)

// Might need to trim <cr> from responses
// HELPER FUNCTIONS
func convertAndSend(socketKey string, message string) bool {
	// Attero Tech devices require a <CR> at the end of commands
	message += "\r"

	sent := framework.WriteLineToSocket(socketKey, message)

	return sent
}

func readAndConvert(socketKey string) (string, error) {
	response := framework.ReadLineFromSocket(socketKey)

	response = strings.Trim(response, "\r\n\x00")

	if response == "" {
		errMsg := "45h3dr - Response was blank "
		framework.AddToErrors(socketKey, errMsg)
		return "unknown", errors.New(errMsg)
	}

	return response, nil
}

// SET FUNCTIONS
// Clear all bluetooth pairings
func setClear(socketKey string, state string) (string, error) {
	function := "setClear"

	value := "notok"
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		if state == `"true"` {
			value, err = setClearDo(socketKey)
		} else {
			value = "ok"
		}
		if value != "ok" { // Something went wrong - perhaps try again
			framework.Log(function + " - 356fds retrying clear operation")
			maxRetries--
			time.Sleep(1 * time.Second)
			if maxRetries == 0 {
				errMsg := fmt.Sprintf(function + "6fggh2sa - max retries reached")
				framework.AddToErrors(socketKey, errMsg)
			}
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func setClearDo(socketKey string) (string, error) {
	function := "setClearDo"
	message := "CBC"

	sent := convertAndSend(socketKey, message)

	if sent != true {
		errMsg := fmt.Sprintf(function + " - b54hgz - error sending command to clear bluetooth pairings")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	response, err := readAndConvert(socketKey)

	if err != nil {
		errMsg := fmt.Sprintf(function + " - hsfdt34 - error reading response: " + err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	framework.Log(function + " - Response: " + fmt.Sprint(response))

	// If we got here, the response was good, so successful return with the state indication
	return "ok", nil
}

// Start bluetooth pairing mode
func setPairing(socketKey string, state string) (string, error) {
	function := "setPairing"

	value := "notok"
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		if state == `"true"` {
			value, err = setPairingDo(socketKey)
		} else {
			value, err = setDisconnectDo(socketKey)
		}
		if value != "ok" { // Something went wrong - perhaps try again
			framework.Log(function + " - 67df236 retrying set pairing operation")
			maxRetries--
			time.Sleep(1 * time.Second)
			if maxRetries == 0 {
				errMsg := fmt.Sprintf(function + "j76jcs - max retries reached")
				framework.AddToErrors(socketKey, errMsg)
			}
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func setPairingDo(socketKey string) (string, error) {
	function := "setPairingDo"
	message := "BTB"

	sent := convertAndSend(socketKey, message)

	if sent != true {
		errMsg := fmt.Sprintf(function + " - 34jckxo - error sending command to set bluetooth pairing")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	response, err := readAndConvert(socketKey)

	if err != nil {
		errMsg := fmt.Sprintf(function + " - 56352dc - error reading response: " + err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	if strings.Contains(response, "NACK") {
		errMsg := "23sdn4 - Command was unsuccessful  "
		framework.AddToErrors(socketKey, errMsg)
		return "unknown", errors.New(errMsg)
	}

	framework.Log(function + " - Response: " + fmt.Sprint(response))

	// If we got here, the response was good, so successful return with the state indication
	return "ok", nil
}

func setDisconnectDo(socketKey string) (string, error) {
	function := "setDisconnectDo"
	message := "BCC"

	sent := convertAndSend(socketKey, message)

	if sent != true {
		errMsg := fmt.Sprintf(function + " - 7c3dsa - error sending command to close bluetooth connection")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	response, err := readAndConvert(socketKey)

	if err != nil {
		errMsg := fmt.Sprintf(function + " - 23x98id - error reading response: " + err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	if strings.Contains(response, "NACK") {
		errMsg := "0dfk2x - Command was unsuccessful  "
		framework.AddToErrors(socketKey, errMsg)
		return "unknown", errors.New(errMsg)
	}

	framework.Log(function + " - Response: " + fmt.Sprint(response))

	// If we got here, the response was good, so successful return with the state indication
	return "ok", nil
}

// GET FUNCTIONS
// Get whether receiver is pairing. If connected or discoverable, true. If idle, false.
func getPairing(socketKey string) (string, error) {
	function := "getPairing"

	value := `"unknown"`
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = getStatusDo(socketKey)
		if value == `"unknown"` { // Something went wrong - perhaps try again
			framework.Log(function + " - 34rd3i retrying get pairing operation")
			maxRetries--
			time.Sleep(1 * time.Second)
			if maxRetries == 0 {
				errMsg := fmt.Sprintf(function + "46gbd5f - max retries reached")
				framework.AddToErrors(socketKey, errMsg)
			}
		} else if value == `"CONNECTED"` || value == `"DISCOVERABLE"` {
			value = `"true"`
			maxRetries = 0
		} else if value == `"IDLE"` {
			value = `"false"`
			maxRetries = 0
		} else {
			maxRetries = 0
		}
	}
	framework.Log("returning pairing value: " + value)
	framework.Log(fmt.Sprintf("Pair Response (raw): [%q]", value))

	return value, err
}

// Get the bluetooth status. Options are Idle, Discoverable, and Connected
func getStatus(socketKey string) (string, error) {
	function := "getStatus"

	value := `"unknown"`
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = getStatusDo(socketKey)
		if value == `"unknown"` { // Something went wrong - perhaps try again
			framework.Log(function + " - 34rd3i retrying status operation")
			maxRetries--
			time.Sleep(1 * time.Second)
			if maxRetries == 0 {
				errMsg := fmt.Sprintf(function + "sh4hd3 - max retries reached")
				framework.AddToErrors(socketKey, errMsg)
			}
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func getStatusDo(socketKey string) (string, error) {
	function := "getStatusDo"
	message := "BTS"

	sent := convertAndSend(socketKey, message)

	if sent != true {
		errMsg := fmt.Sprintf(function + " - 94dk34 - error sending status query")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	response, err := readAndConvert(socketKey)
	if err != nil {
		errMsg := fmt.Sprintf(function + " - 4udj4 - error reading response: " + err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	if strings.Contains(response, "NACK") {
		errMsg := "34jfk - Command was unsuccessful  "
		framework.AddToErrors(socketKey, errMsg)
		return "unknown", errors.New(errMsg)
	}

	value := "unknown"

	// Response will be in the format "ACK BTS 1"
	// Might need to trim <cr>
	responseArray := strings.Split(response, " ")
	statusResponse := responseArray[2]

	if statusResponse == "0" {
		value = "IDLE"
	} else if statusResponse == "1" {
		value = "DISCOVERABLE"
	} else if statusResponse == "2" || statusResponse == "3" || statusResponse == "4" || statusResponse == "5" {
		value = "CONNECTED"
	} else {
		errMsg := function + " - status value was not 0, 1, 2, 3, 4, or 5"
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	// If we got here, the response was good, so successful return with the state indication
	return `"` + value + `"`, nil
}

func getDeviceName(socketKey string) (string, error) {
	function := "getDeviceName"
	message := "BTCDN"

	sent := convertAndSend(socketKey, message)

	if sent != true {
		errMsg := fmt.Sprintf(function + " - 567gsdd - error sending device name query")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	response, err := readAndConvert(socketKey)
	if err != nil {
		errMsg := fmt.Sprintf(function + " - 3634v - error reading response: " + err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	value := "unknown"

	// If the string contains NACK, a device is not connected. Return 'unknown'.
	if !strings.Contains(response, "NACK") {
		// Response will be in the format "ACK BTCDN "Device Name"<CR>"
		responseArray := strings.SplitN(response, "\"", 2)
		if len(responseArray) > 1 {
			deviceNameResponse := responseArray[1]
			if deviceNameResponse == "" {
				value = "unknown"
			} else {
				value = strings.Trim(deviceNameResponse, "\"\r\n")
				value = strings.ReplaceAll(value, "?", "'")
			}
		}
	}

	// If we got here, the response was good, so successful return with the state indication
	return `"` + value + `"`, nil
}

func getMusicInfo(socketKey string) (string, error) {
	function := "getMusicInfo"
	songNameMessage := "BTSONG"
	artistNameMessage := "BTARTIST"
	value := "unknown"
	songName := "unknown"
	artistName := "unknown"

	// Send query for song name
	sent := convertAndSend(socketKey, songNameMessage)

	if sent != true {
		errMsg := fmt.Sprintf(function + " - 345dwdt - error sending song name query")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	response, err := readAndConvert(socketKey)
	if err != nil {
		errMsg := fmt.Sprintf(function + " - oitj5k - error reading response: " + err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	// If the string contains NACK, a device is not connected. Return 'unknown'.
	if !strings.Contains(response, "NACK") {
		// Response will be in the format "ACK BTSONG "Song Name"<CR>"
		responseArray := strings.SplitN(response, "\"", 2)
		if len(responseArray) > 1 {
			songNameResponse := responseArray[1]
			if songNameResponse != "" {
				songName = strings.Trim(songNameResponse, "\"\r\n")
			}
		}
	}
	// Send query for artist name
	sent = convertAndSend(socketKey, artistNameMessage)

	if sent != true {
		errMsg := fmt.Sprintf(function + " - 4ifm56 - error sending artist name query")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	response, err = readAndConvert(socketKey)
	if err != nil {
		errMsg := fmt.Sprintf(function + " - 450ck4 - error reading response: " + err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	// If the string contains NACK, a device is not connected. Return 'unknown'.
	if !strings.Contains(response, "NACK") {
		// Response will be in the format "ACK BTARTIST "Artist Name"<CR>"
		responseArray := strings.SplitN(response, "\"", 2)
		if len(responseArray) > 1 {
			artistNameResponse := responseArray[1]
			if artistNameResponse != "" {
				artistName = strings.Trim(artistNameResponse, "\"\r\n")
			}
		}
	}

	// Either just return the song name or return "songname - artistname"
	if songName != "unknown" {
		value = songName
	}
	if artistName != "unknown" {
		value = value + " - " + artistName
	}

	// If we got here, the response was good, so successful return with the state indication
	return `"` + value + `"`, nil
}
