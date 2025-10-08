一、TT 生成器（Generator）实现要点

1) 输入/输出与错误
	•	输入：rng, ProfileSample{ Vars, MaxDepth, AllowedOps }, plan{ QType, Intent, MCCorrectCount }
	•	输出：CandidatePools{ TruthTablePools{ Formula, Vars, TrueSet, FalseSet } }
	•	可能错误：
	•	ErrDegenerateFormula：恒真/恒假导致池无法满足本次题型/问法
	•	ErrInsufficientCandidates：TrueSet/FalseSet 数量不足以凑齐 4 个选项
	•	ErrGenerationBudgetExceeded：重试超限（防止死循环）

常量建议：MaxRegenAttempts=64（生成公式）、MaxPlanRetries=8（服务层重试次数）。

2) 变量名与顺序
	•	变量名固定表：["p","q","r","s","t"]（最多 5 个）
	•	按采样的 Vars 取前 k 个，保持顺序用于：
	•	枚举赋值时的比特位顺序（稳定性）
	•	赋值字符串渲染时的顺序（稳定、可复现）

3) AST 随机生成（受 AllowedOps / MaxDepth 约束）
	•	递归生成函数 gen(depth)：
	•	终止条件：depth==0 或以小概率“提早停”（例如 25%）→ 生成 Var
	•	运算符采样：从 AllowedOps 均匀抽取
	•	若是 Not：生成 Not(gen(depth-1))
	•	若是二元：生成 op(gen(depth-1), gen(depth-1))
	•	变量覆盖：生成后统计 VarsUsed
	•	若 |VarsUsed|<2 且 Vars≥2，重新生成（提升题目信息量）

说明：不做任何化简（如分配律）以免引入复杂度；统一由字符串化来保证展示清晰。

4) 真值表计算与池构建
	•	枚举 2^Vars 个赋值（mask 从 0 到 2^k-1）：
	•	赋值映射：assign[var[i]] = (mask>>i)&1 == 1
	•	Eval(Formula, assign) 得到布尔值
	•	渲染赋值为字符串："p=T, q=F, r=T"（变量按固定顺序；T/F 大写）
	•	放入 TrueSet 或 FalseSet
	•	退化判定：
	•	len(TrueSet)==0 或 len(FalseSet)==0 → 视为恒假/恒真 → 重新生成公式

5) 可行性检查（按当前 plan 的 Intent/Type）
	•	你无需在生成器内处理所有问法，但最稳的做法是：按 plan 检查可行性，不够就重生
	•	SC + TT_TRUE_ASSIGNMENTS：需要 TrueSet≥1 && FalseSet≥3
	•	SC + TT_FALSE_ASSIGNMENTS：需要 FalseSet≥1 && TrueSet≥3
	•	MC(k) + TT_TRUE_ASSIGNMENTS：需要 TrueSet≥k && FalseSet≥4-k
	•	MC(k) + TT_FALSE_ASSIGNMENTS：需要 FalseSet≥k && TrueSet≥4-k
	•	TF（TT_EVAL_AT_ASSIGNMENT）：无数量要求（任意公式都可）
	•	不满足 → continue 重生，直到成功或超出 MaxRegenAttempts

6) Blueprint hints（建议存的元信息）
	•	vars: 变量名数组（顺序）
	•	truth_table_bits: 把真值表编码为位串（按固定顺序的掩码从 0→2^k-1；真=1，假=0）
	•	true_count / false_count
	•	（可选）formula_str：公式的字符串化结果

⸻

二、字符串化规范（AST → 可读字符串）

为保证稳定性与去重，采用一致的括号策略：
	•	一元：¬ + 子式（子式若为 Var 直接拼，若为复合式无须额外括号）
	•	二元：一律加括号：(left ∘ right)；运算符集合 ∧ ∨ → ↔
	•	Var：变量名（p/q/r/...）

例：(p ∧ (¬q)) → (r ∨ p) 会渲染成 ((p ∧ ¬q) → (r ∨ p))（双括号无妨，关键是一致）。

⸻

三、Formatter 的 TT 逻辑（固定装配）

1) TF（TT_EVAL_AT_ASSIGNMENT）
	•	选项恒为：["True","False"]
	•	α 的选择：从全体赋值均匀抽一条（渲染为 alpha 字符串填写到模板）
	•	正确索引：计算 Eval(Formula, alpha) 为真 → 索引 0，否则索引 1
	•	dataForTemplate：{ "F": formula_str, "alpha": alpha_str }

2) SC
	•	判定方向取决于 Intent：
	•	TT_TRUE_ASSIGNMENTS：从 TrueSet 抽 1，FalseSet 抽 3
	•	TT_FALSE_ASSIGNMENTS：从 FalseSet 抽 1，TrueSet 抽 3
	•	随机无放回抽取；合并后整体打乱；记录正确项的索引（唯一）
	•	dataForTemplate：{ "F": formula_str }（赋值写在选项里）

3) MC（k∈{2,3}，由 planner 抽）
	•	TT_TRUE_ASSIGNMENTS：TrueSet 抽 k，FalseSet 抽 4-k
	•	TT_FALSE_ASSIGNMENTS：FalseSet 抽 k，TrueSet 抽 4-k
	•	打乱；记录所有正确索引（长度 k）
	•	dataForTemplate：同 SC

去重：TrueSet/FalseSet 本身已无重复；装配时不需额外去重。

⸻

四、Validator 在 TT 中的职责
	•	Eval：按 AST 正确求值（我们后面会在 Equivalence/Inference 再用）
	•	ValidateQuestion：
	•	选项数量和正确索引数量与题型一致
	•	索引范围合法
	•	选项无重复（保险起见）
	•	强语义校验（TT 专属）：
	•	TF：CorrectIndex 与 Eval(F, α) 一致
	•	SC/MC：逐个选项解析赋值字符串 → Eval(F, assign) 与“被标记为正确”一致（如不一致，直接报错）

⸻

五、Orchestrator 中与 TT 相关的细节
	•	ProfileSample 采样：
	•	Vars := sample(config.vars_dist)
	•	MaxDepth := sample(config.max_depth_dist)
	•	AllowedOps := parseOps(config.allowed_ops)（映射到 NodeKind）
	•	调 TT 生成器后，构造 dataForTemplate：
	•	F：formula_str（统一字符串化函数）
	•	alpha（仅 TF）：由 Formatter 内部抽到后回写 Blueprint（或提前在 service 中抽，再传给 Formatter；两种都可以）
	•	Blueprint.PoolsSummary 可记录：true_count/false_count，以及如果 TF 用到 α，则把 alpha 一并存入。

⸻

六、单测清单（TT 部分）
	1.	AST/Eval 基本正确性
	•	¬p 在 p=T/F 下返回 F/T
	•	(p → q) 的真值表与定义一致
	•	(p ↔ q) 的真值表与定义一致
	2.	字符串化稳定性
	•	对同一 AST 多次渲染得到相同字符串
	•	常见结构的渲染带括号（避免歧义）
	3.	生成器可行性
	•	在 Easy/Medium/Hard 下，连续 1000 次生成：
	•	不出现 ErrGenerationBudgetExceeded（概率极低）
	•	输出 TrueSet/FalseSet 均非空
	•	在指定 plan=SC/MC(k) 和两种 TT intent 下，能大概率一次成功（统计成功率）
	4.	Formatter 正确性
	•	TF：随机 100 次，CorrectIndex 与 Eval 一致
	•	SC/MC：随机 100 次，每个选项的真假与 CorrectAnswerIndex 集合一致；数量与题型一致
	5.	端到端最小闭环
	•	固定 seed → GenerateQuestion(Easy, TT_TRUE_ASSIGNMENTS, SC) 返回同一题干/选项（可复现）

⸻

七、失败与回退（TT 特化）
	•	优先在 生成器内部重生公式（最多 MaxRegenAttempts 次）
	•	若仍然失败（特别是 MC=3 在低变量情况下偶尔会“真的凑不齐”）：
	•	Orchestrator 将 MC=3 → MC=2 再来一次
	•	仍失败 → 切换 Type=SC（同 Intent）或按 planner 权重换 Intent（如从 TRUE_ASSIGNMENTS 换 FALSE_ASSIGNMENTS）
	•	再失败 → 换 Category（但我们现在只开 TT，就认为报错上抛即可）

⸻

八、落地顺序建议（TT）
	1.	写好 AST.Eval 与 AST.String（字符串化）
	2.	实现 TruthTableGenerator.Generate（含重生循环与可行性检查）
	3.	实现 Formatter 的 TF/SC/MC 固定装配逻辑（TT 分支）
	4.	在 Validator 里补 TT 的强语义校验
	5.	在 Service 里打通 “Plan→ProfileSample→Generate→Build→Validate→Blueprint” 流水
	6.	跑通单测清单（至少前 4 项）
