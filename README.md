# singlegame —— 单机化开发主线

> 本目录是**唯一的开发副本**。`../trunk/` 是只读基线，**任何情况下都不修改**。

## 目录

```
singlegame/
├── JE_WLY_Clinet/        Unity 工程（主战场，源自 trunk 的整份拷贝）
├── server/               对照服务端（只读运行，不开发）
├── docs/                 约束 / 决策记录 / 改动登记
└── README.md             本文件
```

## 十条铁律

> 违反任何一条都会出真问题。每条后面写的"症状"是它坏了之后你会看到的现象。

| # | 铁律 | 违反了会怎样 |
|---|---|---|
| 1 | **`../trunk/` 只读**：不改、不拷回、不在里面开 Unity Play | singlegame 的登录态/存档目录被覆盖；改动无据可查 |
| 2 | 上行协议**只允许在两处拦截**：`TCPClientMgr.SendMessage`（TCP）、`HttpRequestHelper`（HTTP） | 散落各界的 `if (IsLocalMode)`，改一处漏十处 |
| 3 | 推送**只走 `UniGameEvent.Send(具体类型)`**，禁止 `JEvent.defaultEvent.Post` | 编译通过、不报错，但界面一片空白（trunk 的 ProtoC 只监听 UniGameEvent） |
| 4 | **不删 UI 模块、不删 prefab**，只隐藏入口 + 空实现 | 编译链崩、prefab 引用断、`Missing` 一片 |
| 5 | **不直接改原配置表**，改动写进 `.../TextAssets/single/` 覆盖层 | 原版与我方永久混杂，以后没法对比也没法回滚 |
| 6 | 动任何**原文件前先登记** `docs/CHANGES.md` | 半年后没人知道这行为什么是这样 |
| 7 | LocalServer 只写**引擎无关纯 C#**；四层分层（Adapter/Service/Repository/Config）；Adapter 层不许有业务计算 | ILRuntime 下莫名报错；上帝类改一处崩一片 |
| 8 | `.meta` **必须入库**；`Library/` `Bundles/` `Hidden~/` `*.bytes` `obj/` **必须忽略** | 资源引用全断（变紫）/ 仓库 3 GB 卡死 |
| 9 | 跑 batchmode 前**关掉 Unity**；同一时刻**只开一个 Unity** | Library 被锁、测试跑出假红 |
| 10 | 每个工作包收尾必须**能进大厅**，并打 tag | 长期处于"编译过但跑不起来"，出问题无法二分 |

## 命令速查

```bash
# 编热更 DLL（改代码后必跑；不用开 Unity）
cd /d I:\AINewWork\maochao_game\game_machao_new\singlegame
BuildHotUpdate_Local.bat                 # 正式
BuildHotUpdate_Local.bat --test          # 带 LOCAL_TEST，跑测试用

# 起 / 停对照服务端（redis + mock + 游戏服）
server\start_server.bat
server\stop_server.bat

# 跑自动化测试（会占用 Unity，先关 GUI）
powershell -ExecutionPolicy Bypass -File .\LocalTestHarness\runner\Invoke-Batchmode.ps1 -Scenario all
powershell -ExecutionPolicy Bypass -File .\LocalTestHarness\runner\Invoke-Batchmode.ps1 -Scenario equip
```

## 两条顺序约束（顺序错了会白折腾）

```
① 先把 trunk/client/BuildHotUpdate_Local.bat 拷到 singlegame/ 根下（不是 singlegame/client/）
② bat 依赖 JE_WLY_Clinet\Library\ScriptAssemblies\*.dll  →  现在已在位，可直接编
③ 但如果哪天 Library 被删了：必须先 Unity 打开编译一次，再跑 bat
```

## 文档分工

| 文件 | 管什么 | 什么时候更新 |
|---|---|---|
| `README.md`（本文件） | 铁律 + 命令 + 目录。**永远保持一屏** | 铁律变动时 |
| `docs/开发约束.md` | 铁律的**原因、强制手段、自检方法** | 同上 |
| `docs/CHANGES.md` | 每个原文件改动一行：文件 / 原因 / 可回滚 | **动原文件时** |
| `docs/adr/ADR-00X.md` | 结构性决策：背景 / 备选 / 决定 / 后果 | 做决策时 |
| `docs/proto-map.md` | CmdId → 服务端 go 实现位置（S1 产出） | 每次翻译协议前查 |

## 当前状态

- [x] 客户端已拷贝到 `JE_WLY_Clinet/`（Library 已按清单清理，保留 Artifacts/PackageCache/AtlasCache/ShaderCache/ScriptAssemblies）
- [x] server 已放入 `server/`
- [ ] `BuildHotUpdate_Local.bat` 尚未拷到本目录
- [ ] git 尚未初始化（见 `docs/开发约束.md` §4）
- [ ] 首次 Unity 打开 + 6 条冒烟验收
