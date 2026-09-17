# 基线存档（改造前状态）

**Feature**: specs/001-maven-module-deps | **Date**: 2026-09-17 | **任务**: T001

US1 的 RED 证据：改造前"一条命令构建全部后端模块"不可达。

## 形态一：根构建只含 1 个项目（根 POM 非聚合器）

```text
$ ./mvnw validate   （仓库根）
[INFO] Scanning for projects...
[INFO] ------------------< org.juling.ecboot:ecboot-parent >-------------------
[INFO] Building ecboot-parent 0.0.1-SNAPSHOT
[INFO]   from pom.xml
[INFO] --------------------------------[ pom ]---------------------------------
[INFO] BUILD SUCCESS
[INFO] Total time:  0.123 s
```

仅 ecboot-parent 自身参与构建——10 个子模块对 reactor 不可见。

## 形态二：子模块父 POM 不可达（`<relativePath/>` 阻断解析）

```text
$ cd dependencies && ../mvnw validate
[ERROR] The build could not read 1 project -> [Help 1]
[ERROR]   The project org.juling.ecboot:ecboot-dependencies:0.0.1-SNAPSHOT
[ERROR]   (D:\code\git\ecboot2\dependencies\pom.xml) has 1 error
[ERROR]     Non-resolvable parent POM for org.juling.ecboot:ecboot-dependencies:
[ERROR]     0.0.1-SNAPSHOT: The following artifacts could not be resolved:
[ERROR]     org.juling.ecboot:ecboot-parent:pom:0.0.1-SNAPSHOT (absent):
[ERROR]     Could not find artifact ... and 'parent.relativePath' points at no local POM
[ERROR]     @ line 5, column 10
```

全部 `apps/`、`services/`、`infrastructure/`、`dependencies/` 模块同此失败形态。
