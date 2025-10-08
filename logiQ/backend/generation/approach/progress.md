1. 打基础：公共工具和模型

先把核心的 AST 操作补齐：Eval、Stringify、Equivalent、Derivable 等。Validator 里至少要能对 TT/EQ/INF 做语义判定。
写好“分布采样”等小工具（如果还没实现 SampleWeighted 等），它们会被 Planner/Generator 复用。


2. 实现单个 Generator（建议按 TT → EQ → INF 顺序）

从 Truth Table 开始：根据 ProfileSample 生成公式、构建 TrueSet/FalseSet、做数量校验。跑一两个小测试，确认能产出合法池子。
EQ / INF 可以先写接口桩，后续逐步填充逻辑。
3. Formatter

写装配逻辑，把 CandidatePools + IntentSpec + Prompt 变成 models.Question。确保 SC/MC/TF 的正确索引、选项数量等规则都覆盖，便于 Validator 校验。
4. PromptBuilder / IntentLibrary

已经抽象出来了，确认它们能正确拿模板、替换占位符，产出稳定的 Prompt。
5. Planner & Config Loader

写 Planner 时先调通“读 YAML + SamplePlan”的流程。别忘了在载入后把 PromptBuilder、IntentLibrary 等依赖准备好，方便 Service 直接注入。
6. Service (Orchestrator)

前面模块都稳定后再写 Service 的“串联流程”和回退策略。这样你在 Service 里只需要调用已存在的组件，不会边写边想细节。
先打通最短路径（比如只生成 TT + SC），确认能返回一题；再逐步把 MC 回退、类别切换等逻辑补上。
7. 验证

自己写几个简单的单元/集成测试，固定 seed 下看看产出的题干、选项是否稳定。
如果 Validator 还没实现，暂时在 Service 里加注释或 TODO，提醒最后一定要调用。