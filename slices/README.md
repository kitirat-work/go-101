# Slices

- การประกาศนี้สร้าง slice ของ int โดยไม่กำหนดค่าเริ่มต้นให้กับ x ดังนั้น x จะถูกกำหนดค่าเริ่มต้น (zero value) ของ slice ซึ่งคือ nil

```go
var x []int
fmt.Println(x == nil) // ผลลัพธ์: true
```

- ไม่สามารถใช้ == เพื่อเปรียบเทียบสอง slices ได้ จะเกิดข้อผิดพลาดในเวลาคอมไพล์

```go
x := []int{1, 2, 3, 4, 5}
y := []int{1, 2, 3, 4, 5}
fmt.Println(x == y) // ข้อผิดพลาดในการคอมไพล์
fmt.Println(x == nil) // ผลลัพธ์: false
```

- ตั้งแต่ Go 1.21 มีฟังก์ชันในแพ็กเกจ slices ที่ใช้สำหรับเปรียบเทียบ slices:
    - slices.Equal: เปรียบเทียบสอง slices และคืนค่า true หาก slices มีความยาวเท่ากันและทุกองค์ประกอบเท่ากัน
    
    ```go
    import (
        "fmt"
        "golang.org/x/exp/slices"
    )
    
    x := []int{1, 2, 3, 4, 5}
    y := []int{1, 2, 3, 4, 5}
    z := []int{1, 2, 3, 4, 5, 6}
    
    fmt.Println(slices.Equal(x, y)) // ผลลัพธ์: true
    fmt.Println(slices.Equal(x, z)) // ผลลัพธ์: false
    ```
    
    - slices.EqualFunc: ให้คุณส่งฟังก์ชันเพื่อตัดสินความเท่าเทียม และไม่จำเป็นต้องมีองค์ประกอบของ slice ที่สามารถเปรียบเทียบได้
    
    ```go
    import (
        "fmt"
        "golang.org/x/exp/slices"
    )
    
    x := []int{1, 2, 3, 4, 5}
    y := []int{1, 2, 3, 4, 5}
    
    equalFunc := func(a, b int) bool {
        return a == b
    }
    
    fmt.Println(slices.EqualFunc(x, y, equalFunc)) // ผลลัพธ์: true
    ```
    
- ในอดีตมีการใช้ฟังก์ชัน reflect.DeepEqual จากแพ็กเกจ reflect แต่ในปัจจุบันมีฟังก์ชันที่ดีกว่าในแพ็กเกจ slices

```go
import (
    "fmt"
    "reflect"
)

x := []int{1, 2, 3, 4, 5}
y := []int{1, 2, 3, 4, 5}

fmt.Println(reflect.DeepEqual(x, y)) // ผลลัพธ์: true
```

- **การใช้ฟังก์ชัน len กับ Slices ในภาษา Go**
    - หากส่ง slice ที่เป็น nil จะคืนค่า 0
- ฟังก์ชัน len ช่วยให้นักพัฒนาสามารถหาความยาวของ array, slice, string, map, และจำนวนค่าที่รอใน channel ได้อย่างสะดวก
- การใช้งาน len กับชนิดข้อมูลที่ไม่รองรับจะทำให้เกิดข้อผิดพลาดในเวลาคอมไพล์
- **ความจุ (Capacity) ของ Slices ในภาษา Go**
    - การใช้ฟังก์ชัน append เพื่อเพิ่มค่าใน slice จะเพิ่มความยาวของ slice ขึ้นตามจำนวนค่าที่เพิ่ม
    - เมื่อความยาวของ slice ถึงความจุสูงสุดแล้ว หากเพิ่มค่าลงใน slice เพิ่มเติม ฟังก์ชัน append จะใช้ runtime ของ Go เพื่อจัดสรร array ใหม่ที่มีความจุใหญ่ขึ้น
    - ค่าจาก array เดิมจะถูกคัดลอกไปยัง array ใหม่ และค่าใหม่จะถูกเพิ่มที่ปลายของ array ใหม่
    - slice จะถูกอัพเดตให้ชี้ไปยัง array ใหม่ และ slice ที่อัพเดตแล้วจะถูกคืนค่า
    - Go runtime จะเพิ่มขนาดของ slice มากกว่าหนึ่งในแต่ละครั้งที่ความจุหมด โดยกฎใน Go 1.18 คือ:
        - หากความจุน้อยกว่า 256 จะเพิ่มความจุเป็นสองเท่า
        - หากความจุมากกว่า 256 จะเพิ่มความจุโดยใช้สูตร (current_capacity + 768)/4
- **Go Runtime**
    - **บทบาทของ Go Runtime:**
        - **การจัดสรรหน่วยความจำและการเก็บขยะ (Garbage Collection)**: จัดการการใช้งานและการคืนหน่วยความจำ
        - **การสนับสนุนการทำงานพร้อมกัน (Concurrency Support)**: ช่วยให้การทำงานหลายๆ อย่างในเวลาเดียวกันง่ายขึ้น
        - **การสนับสนุนเครือข่าย (Networking)**: ให้ฟังก์ชันในการทำงานกับเครือข่าย
        - **การนำเข้าและใช้งานประเภทข้อมูลและฟังก์ชันในตัว (Built-in Types and Functions)**
    - **ลักษณะพิเศษของ Go Runtime:**
        - Go runtime ถูกคอมไพล์รวมเข้ากับทุก binary ของ Go ซึ่งแตกต่างจากภาษาที่ใช้ virtual machine ที่ต้องติดตั้งแยกต่างหากเพื่อให้โปรแกรมทำงานได้
        - **ข้อเสีย = binary ของโปรแกรม Go แม้จะเป็นโปรแกรมที่ง่ายที่สุด ก็มีขนาดประมาณ 2 MB เนื่องจากการรวม runtime เข้าไป**
- **การใช้ฟังก์ชัน make ในภาษา Go**
    
    ```go
    x := make([]int, 5) // len 5 cap 5 x = [0,0,0,0,0]
    x := make([]int, 5, 10) // len 5 cap 10 x = [0,0,0,0,0] ถ้าเอาไป append ต้องรู้ว่ามี 0 นำหน้านะ
    x := make([]int, 0, 10) // len 0 cap 10 x = []
    ```
    
    - **อย่ากำหนดความจุ (Capacity) น้อยกว่าความยาว (Length)**:
        - หากกำหนดความจุน้อยกว่าความยาวด้วยค่าคงที่หรือลิตเตอรัลตัวเลข จะเกิดข้อผิดพลาดในเวลาคอมไพล์
        - หากกำหนดความจุน้อยกว่าความยาวด้วยตัวแปร โปรแกรมจะ panic ในเวลา runtime
- **การล้างค่าใน Slice ด้วยฟังก์ชัน clear**
    - ทำทำไม ฮ่าๆ
    
    ```go
    package main
    
    import (
        "fmt"
    )
    
    func main() {
        s := []string{"first", "second", "third"}
        fmt.Println(s, len(s)) // ผลลัพธ์: [first second third] 3
    
        clear(s)
        fmt.Println(s, len(s)) // ผลลัพธ์: ["", "", ""] 3
    }
    ```
    
- **การสร้าง Slice จาก Slice**
    - เมื่อคุณสร้าง slice จาก slice เดิม, slice ทั้งสองจะแชร์หน่วยความจำกัน นั่นหมายความว่าการเปลี่ยนแปลงค่าของ slice หนึ่งจะส่งผลกระทบต่ออีก slice หนึ่งที่แชร์หน่วยความจำ

```go
x := []string{"a", "b", "c", "d"}
y := x[:2]
z := x[1:]
x[1] = "y"
y[0] = "x"
z[1] = "z"

fmt.Println("x:", x)
fmt.Println("y:", y)
fmt.Println("z:", z)
//x: [x y z d]
//y: [x y]
//z: [y z d]

//append
x := []string{"a", "b", "c", "d"}
y := x[:2]
fmt.Println(cap(x), cap(y)) // ผลลัพธ์: 4 4

y = append(y, "z")
fmt.Println("x:", x)
fmt.Println("y:", y)
/*
4 4
x: [a b z d]
y: [a b z]
*/
```

- **ใช้ full slice expression เพื่อป้องกันปัญหาจาก append**
    - เพื่อป้องกันปัญหาจากการใช้ append กับ slice ที่ถูกสร้างจาก slice อื่น ควรใช้ full slice expression ซึ่งประกอบด้วยสามส่วน: จุดเริ่มต้น, จุดสิ้นสุด, และความจุสูงสุดที่ต้องการใช้
    
    ```go
    x := make([]string, 0, 5)
    x = append(x, "a", "b", "c", "d")
    y := x[:2:2]
    z := x[2:4:4]
    
    y = append(y, "I", "j", "k")
    x = append(x, "x")
    z = append(z, "y")
    
    fmt.Println("x:", x)
    fmt.Println("y:", y)
    fmt.Println("z:", z)
    /*
    x: [a b c d x]
    y: [a b I j k]
    z: [c d y]
    */
    ```
    
- **การใช้ฟังก์ชัน copy ในภาษา Go**
    - จะเป็นการลอก value มา ไม่แชร์ mem
    
    ```go
    package main
    
    import (
        "fmt"
    )
    
    func main() {
        x := []int{1, 2, 3, 4}
        d := [4]int{5, 6, 7, 8}
        y := make([]int, 2)
    
        // คัดลอกค่าจาก array ไปยัง slice
        copy(y, d[:])
        fmt.Println(y) // ผลลัพธ์: [5 6]
    
        // คัดลอกค่าจาก slice ไปยัง array
        copy(d[:], x)
        fmt.Println(d) // ผลลัพธ์: [1 2 3 4]
    }
    ```
    
- Converting Arrays to Slices
    - การแปลง array เป็น slice มีคุณสมบัติการแชร์หน่วยความจำเช่นเดียวกับการแปลง slice เป็น slice ดังนั้น การเปลี่ยนแปลงค่าของ array จะส่งผลต่อ slice ที่แชร์หน่วยความจำกัน
    
    ```go
    x := [4]int{5, 6, 7, 8}
    y := x[:2]   // แปลงสองค่าแรก
    z := x[2:]   // แปลงสองค่าหลัง
    --------------------------
    package main
    
    import (
        "fmt"
    )
    
    func main() {
        x := [4]int{5, 6, 7, 8}
        y := x[:2]
        z := x[2:]
    
        x[0] = 10
    
        fmt.Println("x:", x) // ผลลัพธ์: [10 6 7 8]
        fmt.Println("y:", y) // ผลลัพธ์: [10 6]
        fmt.Println("z:", z) // ผลลัพธ์: [7 8]
    }
    ```
    
- Converting Slices to Arrays
    - เมื่อแปลง slice เป็น array ข้อมูลใน slice จะถูกคัดลอกไปยังหน่วยความจำใหม่ ซึ่งหมายความว่าการเปลี่ยนแปลงใน slice จะไม่ส่งผลกระทบต่อ array และในทางกลับกัน
    - ขนาดของ array ต้องถูกระบุในเวลาคอมไพล์
    - ขนาดของ array ที่สร้างจาก slice ต้องไม่ใหญ่กว่าความยาวของ slice มิฉะนั้นจะเกิดข้อผิดพลาดในเวลารันไทม์
    
    ```go
    package main
    
    import (
        "fmt"
    )
    
    func main() {
        xSlice := []int{1, 2, 3, 4}
        xArray := [4]int(xSlice)
        smallArray := [2]int{xSlice[0], xSlice[1]}
        
        xSlice[0] = 10
        
        fmt.Println(xSlice)     // ผลลัพธ์: [10 2 3 4]
        fmt.Println(xArray)     // ผลลัพธ์: [1 2 3 4]
        fmt.Println(smallArray) // ผลลัพธ์: [1 2]
    }
    ```
    
    - สามารถแปลง slice เป็น pointer ของ array ได้โดยใช้ type conversion
    - การแปลงนี้ทำให้หน่วยความจำถูกแชร์กันระหว่าง slice และ pointer ของ array ดังนั้นการเปลี่ยนแปลงค่าในหนึ่งจะส่งผลต่ออีกหนึ่ง