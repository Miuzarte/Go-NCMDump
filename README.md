# Go-NCMDump

## At first:

仓库不提供解密 NCM 所需要的 `CoreKey` `MetaKey`, 这俩用搜索引擎都很好找

## Usage

### 直接运行: 

会遍历 `config.inputDir` 目录下的所有 `.ncm` 文件并输出到 `config.outputDir`, `config.inputDir` 为空则使用可执行文件所在的目录(不是工作目录), 导出文件名重复时会自动跳过

### Go-NCMDump "xxx.ncm" "yyy.ncm" "NCMzzz":

会直接输出到对应输入文件的相同位置, 存在文件夹则遍历

### 拖拽文件或文件夹到可执行文件上

同上

## Config

### config.coverOutput

同时输出专辑封面到 `config.outputDir`

### config.coverEmbed

通过 $PATH 调用 `ffmpeg` 将封面嵌入到输出的媒体

### config.highDefinitionCover

将提取出的封面替换为在线拉取的高清版本

### config.multiThread

顾名思义