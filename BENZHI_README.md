基于 Go 实现的风电机组变桨偏航与并网控制平台项目，一款后端服务，完成变桨、偏航、液压制动、测风、变流器并网、塔筒振动联锁、解缆保护与运行审计的 JSON API 控制台。
基于 Go 实现的风电机组变桨偏航与并网控制平台项目，一款后端服务，完成变桨、偏航、液压制动、测风、变流器并网、塔筒振动联锁、解缆保护与运行审计的 JSON API 控制台。

## 项目简介

WindTurbineCtl 面向并网型风电机组的主控逻辑，围绕「测风 → 变桨 → 转速 → 并网 → 发电 → 脱网顺桨」的控制链，把桨叶、偏航、液压制动、测风装置、变流器、塔筒振动、解缆保护和电网侧状态组织成一组可独立校验的组件。控制台只提供 JSON API，不包含前端页面，所有运行状态与操作记录通过 HTTP 接口读写，文件型存储负责把关键状态落盘以便重启恢复。

## 功能模块

- 变桨控制：桨距角调节、顺桨动作、最小桨距到达判定、启动联锁复位。
- 偏航控制：机舱对风跟踪、手动偏航指令仲裁、电缆缠绕计数与解缆回绕。
- 液压制动：制动压力落盘、蓄能器充压、制动闩锁的投入与恢复解除。
- 测风与功率：风速风向采样、功率曲线插值、限功率目标计算与曲线重标定。
- 变流器并网：切入合闸顺序、脱网顺桨顺序、断路器开合状态管理。
- 塔筒与电网：振动联锁按闩锁管理，电网失电触发保护顺序。
- 运行审计：全部控制动作写入审计记录，按序列号查询最近事件。

## 构建

项目使用 vendor 离线依赖，构建前无需访问网络：

```text
go build -mod=vendor ./...
go test -mod=vendor ./...
go vet -mod=vendor ./...
```

## 运行

```text
go run -mod=vendor ./cmd/wtc
```

默认监听 8080 端口，数据目录为 ./data。可用环境变量覆盖：

- WTC_ADDR：监听地址，默认 :8080。
- WTC_DATA_DIR：数据目录，默认 data。
- WTC_SAMPLE_WINDOW：测风采样窗口，默认 10。
- WTC_AUDIT_BUFFER：审计记录保留条数，默认 200。

## API 一览

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | /api/status | 整机运行状态 |
| GET | /api/audit | 最近审计事件 |
| POST | /api/pitch/move | 变桨移动 |
| POST | /api/pitch/demand | 变桨指令（先充压） |
| POST | /api/pitch/feather | 顺桨 |
| POST | /api/pitch/reset | 启动联锁复位 |
| POST | /api/yaw/turn | 偏航转动 |
| POST | /api/yaw/untwist | 解缆回绕 |
| POST | /api/brake/pressure | 设置并落盘制动压力 |
| POST | /api/brake/charge | 蓄能器充压 |
| POST | /api/brake/latch | 制动闩锁投入或解除 |
| POST | /api/wind/sample | 写入测风采样 |
| POST | /api/power/target | 按风速计算功率目标 |
| POST | /api/cutin | 切入并网 |
| POST | /api/breaker/open | 手动断开断路器 |
| POST | /api/limit | 限功率指令 |
| POST | /api/trip | 电网失电保护 |
| POST | /api/recalibrate | 功率曲线重标定 |
| POST | /api/tower/vibration | 塔筒振动采样 |

## 目录结构

```text
cmd/wtc        程序入口、配置与组件装配
internal/ns    命名空间与限额模型
internal/blade 桨叶状态
internal/pitch 变桨控制
internal/yaw   偏航控制
internal/brake 液压制动
internal/wind  测风装置
internal/conv  变流器
internal/tower 塔筒振动
internal/cable 解缆保护
internal/grid  电网侧
internal/store 文件型持久化
internal/audit 运行审计
internal/console JSON API 控制台
```

## 容器镜像

```text
sh build_benzhi_docker.sh
```

镜像基于 golang 基础镜像，内部启用 GOPROXY=off 并使用 vendor 离线构建，可直接以 CMD 启动服务进程。
