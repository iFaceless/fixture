```
  __ _      _
 / _(_)_  _| |_ _   _ _ __ ___
| |_| \ \/ / __| | | | '__/ _ \
|  _| |>  <| |_| |_| | | |  __/
|_| |_/_/\_\\__|\__,_|_|  \___|
```

# 介绍

我们在 Python 项目中写单元测试时，一般需要准备测试数据并在启动单元测试时使用 `@fixture('foo.sql')` 的方式导入到数据库。

让我们回顾下，使用 SQL 格式的 fixture 数据存在什么问题：
1. 造测试数据非常麻烦，如果某个表的字段较多时，更为严重，可能存在某些字段填的值没和期望的列名对应上；
2. 纯 SQL 虽然可以方便地从数据库直接复制导入，但是后期编辑和维护比较麻烦。

受到开源项目 [theiconic/fixtures](https://github.com/theiconic/fixtures) 和 [zhulongcheng/testsql](https://github.com/zhulongcheng/testsql) 的启发，打算使用 YAML/JSON 编写测试数据。当然，这里极力推荐 YAML 格式，方便人肉编辑！

# 功能

**注意**：*目前仅支持 MySQL 或者符合 MySQL 协议的数据库*

- **优雅的接口**：方便在测试代码中使用 fixture 数据
- **支持 JSON/YAML/SQL 格式测试数据导入数据库**
- **支持从数据库指定表中生成 fixture 数据（支持 JSON/YAML/SQL 格式导出）**
- **格式可扩展**：除了默认支持的 `JSON/YAML/SQL` 格式外，也支持自定义格式，只需要实现相关接口即可（直接提 PR）


# 安装

```sh
go get -u github.com/iFaceless/fixture
```

# 示例

首先，假设包含测试数据的项目结构如下：

```
.
├── testdata 这里放测试数据
│   ├── fixtures
│   │   ├── user.yml 文件名默认为表名，每个文件都是针对单个表的测试数据
│   │   ├── topic.json 同样支持 JSON 格式数据
│   │   └── product.yml
│   └── schema.sql
└── thrifts
```

下面我们看在单元测试中如何使用这个工具。

## 配置

首先我们需要告诉该工具一些测试的配置，建议在项目的 `pkg/configs` 包中新建一个文件 `fixture.go`，目录结构类似：

```
configs
└── fixture.go
```

在 `fixture.go` 中添加下面的代码，方便 fixture 启动时知道去哪儿读取 schema 和测试数据配置：

```golang
func NewDefaultTestFixture() *fixture.TestFixture {
	return fixture.New(
		fixture.Database(getTestDBRawURL()),
		fixture.SchemaFilepath(path.Join(getTestDataPath(), "schema.sql")),
		fixture.DataDir(path.Join(getTestDataPath(), "fixtures")),
	)
}

func getTestDBRawURL() string {
	dsnFmt := "mysql://%s:%s@%s/%s?charset=utf8&parseTime=true&loc=Asia/Shanghai"
	return fmt.Sprintf(dsnFmt,
		os.Getenv("MYSQL_USERNAME"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_DATABASE"),
	)
}

// getTestDataPath 返回测试夹具数据路径
func getTestDataPath() string {
	_, f, _, _ := runtime.Caller(1)
	curDir := path.Dir(f)
	// 找到项目 root 目录
	// 类似 ../../curDir
	root := path.Dir(path.Dir(curDir))
	return path.Join(root, "testdata")
}
```

## 使用姿势

```golang
type SuiteExampleTester struct {
	suite.Suite
	tf *fixture.TestFixture
}

func (s *SuiteExampleTester) SetupSuite() {
	s.tf = NewDefaultTestFixture()
}

func (s *SuiteExampleTester) TearDownSuite() {
	s.tf.DropTables()
}

func (s *SuiteExampleTester) TestFoo() {
	scope := s.tf.Use("table_a", "table_b")
	defer scope.Clear()

	// 填写测试逻辑
}

func (s *SuiteExampleTester) TestBar() {
	s.tf.Use("table_a", "table_b").Test(func () {
		// 在该闭包中填写测试逻辑
		// 测试完毕后，会自动清空表
	})
}
```

# 导出

啊，我不想一开始就手写 YAML/JSON/SQL！我希望能利用测试环境中的数据，然后再在输出结果的基础上修改为想要的测试数据！没问题，`fixture` 目前已经支持导出数据为这三种格式啦！

## 安装

```shell
$ go get github.com/iFaceless/fixture/fixturegen
```

检查下是否安装正确（确保你的 $GOPATH 在 $PATH 环境变量中！）

```shell
$ which fixturegen
# 输出结果类似这样：
`$GOPATH`/bin/fixturegen

# 查看下帮助
$ fixturegen --help

Usage of fixturegen:
  -ext string
    	output file extension (e.g. '.yml', '.json'） (default ".yml")
  -o string
    	output directory (default ".")
  -q string
    	custom query sql
  -t string
    	table to be exported
  -url string
    	database connection url
```

## 使用姿势

下面我们看看怎么使用它生成想要的测试数据：

```shell
# 默认会生成 YAML 格式文件，并使用默认的数据查询策略（只查 5 条，并基于 id 排序）
# 输出的文件默认在工作目录下
$ fixturegen -url mysql://localhost:3306/test_todo_api -t user

# 可以指定别的格式输出
$ fixturegen -url mysql://localhost:3306/test_todo_api -t user -ext [.json/.ext/.yml]

# 当然，也可以指定输出数据到别的目录
$ fixturegen -url mysql://localhost:3306/test_todo_api -t user -o path/to/testdata/fixtures

# 哦哦，还可以自定义查询语句（注意：只支持简单的查询语句）
$ fixturegen -url mysql://localhost:3306/test_todo_api -t user -q "SELECT * FROM user WHERE id > 100 LIMIT 10"
```

生成文件示例：

![YAML 文件](http://pho2na29z.bkt.clouddn.com/2018-12-14-21-05-11.png)

![JSON 文件](http://pho2na29z.bkt.clouddn.com/2018-12-14-21-05-28.png)

![SQL 文件](http://pho2na29z.bkt.clouddn.com/2018-12-14-21-05-59.png)

# 主要 API 说明

- `TestFixture.New`: 新建 `TestFixture` 实例，需要用户提供数据库、测试数据配置
- `TestFixture.Use`: 使用指定表的测试数据填充到测试数据库对应表中
- `TestFixture.DropTables`: 用于测试结束后删除测试表（注意，`fixture` 工具不会随意自动删除表，所以作为用户的你需要显式调用才会删除表）
- `TestFixture.TableNames`: 通过 `schema.sql` 读取到的所有表名
- `TestFixture.Config`: 可以获取详细配置信息
- `Scope.Clear`: 用于某个单元测试结束后，清空表数据
- `Scope.Test`: 可接收一个测试函数，运行测试函数后自动清空表

# Help & Dev & Bug Report

如有任何想要的新功能或者是发现的 Bug 均可提到 issue 中~

