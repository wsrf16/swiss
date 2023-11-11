package base62kit

// var base = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
//
// func Encode(number int) string {
//     baseStr := ""
//     for {
//         if number <= 0 {
//             break
//         }
//
//         i := number % 62
//         baseStr = base[i] + baseStr
//         number = (number - i) / 62
//     }
//     return baseStr
// }
//
// func Decode(base62 string) int {
//     rs := 0
//     len := len(base62)
//     f := flip(base)
//     for i := 0; i < len; i++ {
//         rs += f[string(base62[i])] * int(math.Pow(62, float64(len-1-i)))
//     }
//     return rs
// }
//
// func flip(s []string) map[string]int {
//     f := make(map[string]int)
//     for index, value := range s {
//         f[value] = index
//     }
//     return f
// }
