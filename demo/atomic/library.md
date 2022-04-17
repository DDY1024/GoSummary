#### https://pkg.go.dev/sync/atomic
#### 1. 整数相关操作
##### 1.1 整数类型
- `int32`
- `int64`
- `uint32`
- `uint64`

##### 1.2 操作类型
- `Add`
- `Load`
- `Store`
- `CAS`
- `Swap`：直接交换，建议还是 CAS

#### 2. 指针操作
##### 2.1 uintptr
- `Add`
- `Load`
- `CAS`
- `Store`
- `Swap`：直接交换，建议还是 CAS

##### 2.2 unsafe.Pointer
- `Load`
- `CAS`
- `Store`
- `Swap`：直接交换，建议还是 CAS

#### 3. CAS 操作
##### 3.1 整数类型
- `int32`
- `int64`
- `uint32`
- `uint64`

##### 3.2 指针类型
- `unsafe.Pointer`
- `uintptr`

#### 4. atomic.Value 
- `CAS`
- `Load`
- `Store`
- `Swap`