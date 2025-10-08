Equivalence 生成器（无代码、到实现可落地）

0) 目标与接口
	•	输入：rng、ProfileSample{ Vars, MaxDepth, AllowedOps }、plan{ QType, Intent, MCCorrectCount }
	•	输出：CandidatePools{ EquivalencePools{ Target, Vars, EquivPool, NonEquivPool } }
	•	错误：
	•	ErrDegenerateFormula（目标式是恒真/恒假，或导致池凑不齐）
	•	ErrInsufficientCandidates（等价或非等价候选数量不足）
	•	ErrGenerationBudgetExceeded（重生超过上限）

建议：Vars 至少 2（若采样为 1，自动提到 2），避免生成空间过小。

⸻

1) 生成目标式 T（与 TT 复用 AST 生成器）
	•	用 TT 同款 AST 随机器，受 AllowedOps/MaxDepth/Vars 约束。
	•	退化过滤：算一遍真值表，要求 TrueSet>0 && FalseSet>0（避免恒真/恒假）。
	•	变量覆盖：优先确保用到 ≥2 个变量（若可行）。

⸻

2) 候选池设计 & 最小变换/扰动库

2.1 等价候选（EquivPool）：安全变换库（任选 5~8 条足够）

仅做 1 步（默认），少量情况允许 2 步（可配置，但我们当前不做 selection policy，只用到“凑够数量”为止）。
	•	双重否定：¬¬A ≡ A
	•	德摩根：¬(A ∧ B) ≡ (¬A ∨ ¬B)；¬(A ∨ B) ≡ (¬A ∧ ¬B)
	•	蕴含消去：A→B ≡ ¬A∨B
	•	↔ 展开：A↔B ≡ (A→B)∧(B→A)（仅在 AllowedOps 含 IFF 或 T 中已出现 ↔）
	•	交换律：A∧B ≡ B∧A；A∨B ≡ B∨A
	•	结合律：A∧(B∧C) ≡ (A∧B)∧C；A∨(B∨C) ≡ (A∨B)∨C
	•	分配律：A∧(B∨C) ≡ (A∧B)∨(A∧C)；A∨(B∧C) ≡ (A∨B)∧(A∨C)
	•	逆否等价：A→B ≡ ¬B→¬A（若 → 在 AllowedOps 中）

校验：每个变体都做 真值表等价验证；字符串化后去重，不允许与 T 本体完全相同（避免“把 T 当选项”的无效题）。

2.2 非等价候选（NonEquivPool）：单点扰动库（逐一验证“非等价”）
	•	联结词翻转：∧↔∨、→↔↔（把 ↔ 替换成 →，或把 → 换成 ↔）
	•	方向调换：A→B 改为 B→A
	•	否定错位：把 ¬(A ∘ B) 错改成 ¬A ∘ ¬B（∘ 与德摩根相反处置一处）或删掉/加上一个 ¬（只在一个子式上）
	•	变量局部替换：在一个子式内把 p 换成 q（T 若对称，这步可能“意外等价”，故须验证）
	•	括号结构调整：把 (A ∧ (B ∧ C)) 替成 (A ∧ B) ∨ C 等明显改变优先级与运算符的替换
	•	替换一侧子式：如 A↔B 中把 B 替换成 ¬B 或 B∧C

强校验：每个扰动产物都跑 非等价验证（与 T 有一处赋值不同即可），失败就丢弃并换另一扰动。

注意：我们不使用“近错阈值”挑拣；只要“确实非等价”且数量够即可。

⸻

3) 产量与可行性（对接题型/问法）

根据 plan 事先 凑够数量（不做事后 selection）：
	•	SC + EQ_EQUIVALENT：Equiv≥1，NonEquiv≥3
	•	SC + EQ_NONEQUIVALENT：NonEquiv≥1，Equiv≥3
	•	MC(k) + EQ_EQUIVALENT：Equiv≥k，NonEquiv≥4-k（k∈{2,3}）
	•	MC(k) + EQ_NONEQUIVALENT：NonEquiv≥k，Equiv≥4-k
	•	TF(EQ_PAIR_TF)：任取一条等价或非等价的 G 与 T 组成判定对；建议 50/50 随机（也可让 planner 控配比）

不足时：
	1.	先继续生成更多候选（再次变换/扰动）；
	2.	仍不足 → 重生 T；
	3.	多次失败 → 上抛 ErrInsufficientCandidates，由 orchestrator 走回退（MC=3→2，或换 Type/Intent）。

⸻

4) 深度与可读性控制（不使用 selection policy 的简化做法）
	•	变换 首选“等价但不爆炸”的：双重否定、德摩根、蕴含消去、交换/结合；必要时再使用 ↔ 展开/分配。
	•	可设置一个 软上限：产物的 depth ≤ MaxDepth + 1；超过就丢弃重来（不是阈值筛选，只是避免极端展开）。
	•	仍不够则换 T。

⸻

5) 字符串化与去重
	•	采用 统一括号规范（与 TT 一致）：二元一律括号；一元 ¬ 前缀；运算符用 ∧ ∨ → ↔。
	•	用字符串做集合去重；必要时可做 AST 正规化（比如对 ∧/∨ 进行 commutative sort），但非必须。

⸻

6) Blueprint 建议记录
	•	target_str: T 的渲染
	•	equiv_count / nonequiv_count
	•	equiv_variants: 列表（可简要记 rule 与“作用位置”）
	•	nonequiv_variants: 列表（记 perturbation 与位置）
	•	vars: 变量列表与顺序
	•	（可选）depth_of_T 与平均深度

⸻

7) Formatter（等价题）的装配（与 TT 一样固定）
	•	SC：按 intent 方向，从对应池取 1 + 3；整体打乱；记录唯一正确索引
	•	MC(k)：按方向取 k + (4-k)；打乱；记录 k 个正确索引
	•	TF：给出 ["True","False"]；对 (T,G) 调用等价判定，决定正确索引
	•	dataForTemplate：
	•	SC/MC：{"T": target_str}（选项里放候选式字符串）
	•	TF：{"T": target_str, "G": candidate_str}

⸻

8) 单测清单（Equivalence）
	1.	等价验证：对库中每条安全变换，构造 (T, T')，断言 Equivalent(T,T') == true
	2.	非等价验证：对每条扰动，构造 (T, T*)，断言 Equivalent == false（多跑随机）
	3.	可行性：在 Easy/Medium/Hard 下，多次生成，统计失败率（期望极低）；SC、MC(k)、TF 各跑 1000 次
	4.	去重：产生的候选在字符串层面不重复
	5.	不引入非法联结词：产物中 ops ⊆ AllowedOps
	6.	端到端：固定 seed，EQ_EQUIVALENT + SC 返回可复现题干/选项