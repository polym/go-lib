# UPYUN Go SDK

[![Build Status](https://travis-ci.org/upyun/go-sdk.svg?branch=master)](https://travis-ci.org/upyun/go-sdk)

    import "github.com/upyun/go-sdk/upyun"

UPYUN Go SDK, 集成：
- [UPYUN HTTP REST 接口](http://docs.upyun.com/api/rest_api/)
- [UPYUN HTTP FORM 接口](http://docs.upyun.com/api/form_api/)
- [UPYUN 缓存刷新接口](http://docs.upyun.com/api/purge/)
- [UPYUN 分块上传接口](http://docs.upyun.com/api/multipart_upload/)
- [UPYUN 视频处理接口](http://docs.upyun.com/api/av_pretreatment/)

## Projects using this SDK

- [UPYUN Command Tool](https://github.com/polym/upx) by [polym](https://github.com/polym)

## Usage

### 快速上手

```go
```


### 初始化 UpYun

```go
func NewUpYun(config *UpYunConfig) *UpYun
```

`NewUpYun` 初始化 `UpYun`，`UpYun` 是调用又拍云服务的统一入口，`UpYun` 对所有开放的接口都做了支持。

---

### 又拍云 REST API 接口

#### 获取空间存储使用量

```go
func (up *UpYun) Usage() (n int64, err error)
```

#### 创建目录

```go
func (up *UpYun) Mkdir(path string) error
```

#### 上传

```go
func (up *UpYun) Put(config *[PutObjectConfig](#PutObjectConfig)) (err error)
```

#### 下载

```go
func (up *UpYun) Get(config *GetObjectConfig) (fInfo *FileInfo, err error)
```

#### 删除

```go
func (up *UpYun) Delete(config *DeleteObjectConfig) error
```

#### 获取文件信息

```go
func (up *UpYun) GetInfo(path string) (*FileInfo, error)
```

#### 获取文件列表

```go
func (up *UpYun) List(config *GetObjectsConfig) error
```

---

### UPYUN 缓存刷新接口

```go
func (u *UpYun) Purge(urls []string) (string, error)
```

---

### 又拍云表单上传接口

#### 上传文件

```go
func (up *UpYun) FormUpload(config *FormUploadConfig) (*FormUploadResp, error)
```

---

### 又拍云处理接口

#### 提供处理任务

```go
func (up *UpYun) CommitTasks(config *CommitTasksConfig) (taskIds []string, err error)
```

`tasksIds` 是提交任务的编号。通过这个编号，可以查询到处理进度以及处理结果等状态。

#### 获取处理进度

```go
func (up *UpYun) GetProgress(taskIds []string) (result map[string]int, err error)
```

#### 获取处理结果

```go
func (up *UpYun) GetResult(taskIds []string) (result map[string]interface{}, err error)
```

---

### 基本类型

#### UpYun

```go
type UpYunConfig struct {
        Bucket    string                // 云存储服务名（空间名）
        Operator  string                // 操作员
        Password  string                // 密码
        Secret    string                // 表单上传密钥，已经弃用
        Hosts     map[string]string     // 自定义 Hosts 映射关系
        UserAgent string                // HTTP User-Agent 头，默认
}
```

`UpYunConfig` 提供了初始化 `UpYun` 的参数。 需要注意的是，`Secret` 表单密钥已经弃用，如果一定需要使用，需调用 `Use`


#### FileInfo

```go
type FileInfo struct {
        Name        string              // 文件名
        Size        int64               // 文件大小, 目录大小为 0
        ContentType string              // 文件 Content-Type
        IsDir       bool                // 是否为目录
        ETag        string              // ETag 值
        Time        time.Time           // 文件修改时间

        Meta map[string]string          // Metadata 数据

        /* image information */
        ImgType   string
        ImgWidth  int64
        ImgHeight int64
        ImgFrames int64
}
```

#### FormUploadResp

```go
type FormUploadResp struct {
        Code      int      `json:"code"`            // 状态码
        Msg       string   `json:"message"`         // 状态信息
        Url       string   `json:"url"`             // 保存路径
        Timestamp int64    `json:"time"`            // 时间戳
        ImgWidth  int      `json:"image-width"`     // 图片宽度
        ImgHeight int      `json:"image-height"`    // 图片高度
        ImgFrames int      `json:"image-frames"`    // 图片帧数
        ImgType   string   `json:"image-type"`      // 图片类型
        Sign      string   `json:"sign"`            // 签名
        Taskids   []string `json:"task_ids"`        // 异步任务
}
```

`FormUploadResp` 为表单上传的返回内容的格式。其中 `Code` 字段为状态码，可以查看 [API 错误码表](https://docs.upyun.com/api/errno/)

#### PutObjectConfig

```go
type PutObjectConfig struct {
        Path              string                // 云存储中的路径
        LocalPath         string                // 待上传文件在本地文件系统中的路径
        Reader            io.Reader             // 待上传的内容
        Headers           map[string]string     // 请求额外的 HTTP 头
        UseMD5            bool                  // 是否需要 MD5 校验
        UseResumeUpload   bool                  // 是否使用断点续传
        AppendContent     bool                  // 是否是追加文件内容
        ResumePartSize    int64                 // 断点续传块大小
        MaxResumePutTries int                   // 断点续传最大重试次数
}
```

`PutObjectConfig` 提供上传单个文件所需的参数。有几点需要注意:
- `LocalPath` 跟 `Reader` 是一个互斥的关系，如果设置了 `LocalPath`，SDK 就会去读取这个文件，而忽略 `Reader` 中的内容。
- 如果 `Reader` 是一个流／缓冲等的话，需要通过 `Headers` 参数设置 `Content-Length`，SDK 默认会对 `*os.File` 增加该字段。
- [断点续传](https://docs.upyun.com/api/rest_api/#_3)的上传内容必须是 `*os.File`, 断点续传会将文件按照 `ResumePartSize` 进行切割，然后按次序一块一块上传，如果遇到网络问题，会进行重试，重试 `MaxResumePutTries` 次，默认无限重试。
- `AppendContent` 如果是追加文件的话，确保非最后的分片必须为 1M 的整数倍。
- 如果需要 MD5 校验，SDK 对 `*os.File` 会自动计算 MD5 值，其他类型需要自行通过 `Headers` 参数设置 `Content-MD5`


#### GetObjectConfig

```go
type GetObjectConfig struct {
        Path      string                    // 云存储中的路径
        Headers   map[string]string         // 请求额外的 HTTP 头
        LocalPath string                    // 文件本地保存路径
        Writer    io.Writer                 // 保存内容的容器
}
```

`GetObjectConfig` 提供下载单个文件所需的参数。 跟 `PutObjectConfig` 类似，`LocalPath` 跟 `Writer` 是一个互斥的关系，如果设置了 `LocalPath`，SDK 就会把内容写入到这个文件中，而忽略 `Writer`。


#### GetObjectsConfig

```go
type GetObjectsConfig struct {
        Path           string                   // 云存储中的路径
        Headers        map[string]string        // 请求额外的 HTTP 头
        ObjectsChan    chan *FileInfo           // 对象 Channel
        QuitChan       chan bool                // 停止信号
        MaxListObjects int                      // 最大列对象个数
        MaxListTries   int                      // 列目录最大重试次数
        MaxListLevel int                        // 递归最大深度
        DescOrder bool                          // 是否按降序列取，默认为生序

        // Has unexported fields.
}
```

`GetObjectsConfig` 提供列目录所需的参数。当列目录结束后，SDK 会将 `ObjectChan` 关闭掉。


#### DeleteObjectConfig

```go
type DeleteObjectConfig struct {
        Path  string        // 云存储中的路径
        Async bool          // 是否使用异步删除
}
```

`DeleteObjectConfig` 提供删除单个文件／空目录所需的参数。


#### FormUploadConfig

```go
type FormUploadConfig struct {
        LocalPath      string                       // 待上传的文件路径
        SaveKey        string                       // 保存路径
        ExpireAfterSec int64                        // 签名超时时间
        NotifyUrl      string                       // 结果回调地址
        Apps           []map[string]interface{}     // 异步处理任务
        Options        map[string]interface{}       // 更多自定义参数
}
```

`FormUploadConfig` 提供表单上传所需的参数。


#### CommitTasksConfig

```go
type CommitTasksConfig struct {
        AppName   string                            // 异步任务名称
        NotifyUrl string                            // 回调地址
        Tasks     []interface{}                     // 任务数组

        // Naga 相关配置
        Accept    string                            // 回调支持的类型，默认为 json
        Source    string                            // 处理原文件路径
}
```

`CommitTasksConfig` 提供提交异步任务所需的参数。`Accept` 跟 `Source` 仅与异步音视频处理有关。`Tasks` 是一个任务数组，数组中的每一个元素都是任务相关的参数（一般情况下为字典类型）。
