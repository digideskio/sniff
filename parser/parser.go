package parser

import ()

var TLSHeaderLengh = 5

/*
 *
 */
func GetHostname(data []byte) []byte {
	extensions := GetExtensionBlock(data)
	sn := GetSNBlock(extensions)
	sni := GetSNIBlock(sn)
	return sni
}

func GetSNIBlock(data []byte) []byte {
	index := 0

	for {
		if index >= len(data) {
			break
		}
		length := int((data[index] << 8) + data[index+1])
		endIndex := index + 2 + length
		if data[index+2] == 0x00 { /* SNI */
			sni := data[index+3:]
			sniLength := int((sni[0] << 8) + sni[1])
			return sni[2 : sniLength+2]
		}
		index = endIndex
	}
	return []byte{}
}

/* Given an Extensions block, go ahead and find the SN block */
func GetSNBlock(data []byte) []byte {
	index := 0
	extensionLength := int((data[index] << 8) + data[index+1])
	data = data[2:extensionLength]

	for {
		if index >= len(data) {
			break
		}
		length := int((data[index+2] << 8) + data[index+3])
		endIndex := index + 4 + length
		if data[index] == 0x00 && data[index+1] == 0x00 {
			return data[index+4 : endIndex]
		}

		index = endIndex
	}

	/* No one :\ */
	return []byte{}
}

/* */
func GetExtensionBlock(data []byte) []byte {
	/*   data[0]           - content type
	 *   data[1], data[2]  - major/minor version
	 *   data[3], data[4]  - total length
	 *   data[...38+5]     - start of SessionID (length bit)
	 *   data[38+5]        - length of SessionID
	 */
	var index = TLSHeaderLengh + 38

	/* Index is at SessionID Length bit */
	index = index + 1 + int(data[index])

	/* Index is at Cipher List Length bits */
	index = (index + 2 + int((data[index]<<8)+data[index+1]))

	/* Index is now at the compression length bit */
	index = index + 1 + int(data[index])

	/* Now we're at the Extension start */
	return data[index:]
}