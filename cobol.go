/*
Package cobol provides functions for converting text and signed/unsigned integers + decimals to Cobol friendly formats
*/

package cobol

import (
  "log"
  "math/big"
  "strconv"
  "strings"
)

var cacheMap = make(map[string]string)

func pad(input string, paddingString string, length int, direction string) (string) {
  inputLength := len(input)

  if (inputLength > length) {
    return input[0:length]
  }

  var paddingCharacters = strings.Repeat(paddingString, length - inputLength)

  if (direction == "right") {
    return input + paddingCharacters
  }

  return paddingCharacters + input
}

func formatNumericalData(input string, format string, length int) (string) {
  hasDecimal := strings.Contains(format, "V")

  if (hasDecimal) {
    var numDecimalDigits int
    var numIntegralDigits int

    if (format == "9(12)V9(5)") {
      numDecimalDigits = 5
      numIntegralDigits = 12
    } else if (format == "9(15)V99") {
      numDecimalDigits = 2
      numIntegralDigits = 15
    } else if (format == "9(3)V9(6)") {
      numDecimalDigits = 3
      numIntegralDigits = 6
    }

    parsedInput := new(big.Float)
    parsedInput.SetString(input)
    parsedInput.Abs(parsedInput)

    floatString := parsedInput.String()
    splitFloatString := strings.Split(floatString, ".")

    integralPart := splitFloatString[0]
    decimalPart := splitFloatString[1]

    return pad(integralPart, "0", numIntegralDigits, "left") + pad(decimalPart, "0", numDecimalDigits, "right")
  } else {
    parsedInput, _ := strconv.Atoi(input)

    return pad(strconv.Itoa(parsedInput), "0", length, "left")
  }
}

func formatTextData(input string, format string, length int) (string) {
  return pad(input, " ", length, "right")
}

// FormatData takes input, a format like X(50), 9(9), 9(15)V99, and the length of the field from the copybook
// and outputs a formatted value, with memoization
func FormatData(input string, format string, length int) (string) {
  cacheKey := "formatData:" + input + ":" + format + ":" + strconv.Itoa(length)

  if _, ok := cacheMap[cacheKey]; ok {
    return cacheMap[cacheKey]
  }

  isText := format[0] == 'X'
  isNumeric := format[0] == '9'

  if (!isText && !isNumeric) {
    log.Fatal("Invalid format {}\n", format)
  }

  if (isText) {
    cacheMap[cacheKey] = formatTextData(input, format, length)
  } else if (isNumeric) {
    cacheMap[cacheKey] = formatNumericalData(input, format, length)
  }

  return cacheMap[cacheKey]
}

type Cobol interface {
  FormatData(input string, format string, length int) (string)
}
