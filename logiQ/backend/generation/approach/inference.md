inference generation

0）接口目标回顾
	•	输入：rng、ProfileSample{Vars, MaxDepth, AllowedOps, ChainSteps(target)}、plan{QType, Intent, MCCorrectCount}
	•	输出：CandidatePools.Inference{ Premises, Vars, ValidConclusions, InvalidConclusions }
	•	错误：ErrInsufficientCandidates、ErrDegeneratePremises（无满足“前提全真”的赋值）、ErrGenerationBudgetExceeded

统一裁判：Derivable(Premises, C) 按真值表定义：
① 存在至少一个赋值使所有前提为真（非空前提域）；② 在所有“前提全真”的赋值上，结论恒真。

1）最小规则库（模板“积木”）

1.1 有效模板（用于生成 ValidConclusions）
	•	MP（Modus Ponens）：A→B,  A  ⟹  B（链长 1）
	•	MT（Modus Tollens）：A→B,  ¬B  ⟹  ¬A（链长 1）
	•	DS（Disjunctive Syllogism）：A∨B,  ¬A  ⟹  B（链长 1）
	•	HS（Hypothetical Syllogism）：A→B,  B→C  ⟹  A→C（链长 1）
	•	链式（2 步示例）：
	•	A→B, B→C, A  ⟹  C（HS+MP）
	•	A∨B, ¬A, B→C  ⟹  C（DS+MP）
	•	链式（3 步示例）：
	•	A→B, B→C, C→D, A  ⟹  D（MP+MP+MP）
	•	A→B, B→(C∨D), ¬C  ⟹  D（MP+DS 变体）

1.2 常见谬误（用于生成 InvalidConclusions）
	•	Affirming the Consequent：A→B,  B  ⟹  A（无效）
	•	Denying the Antecedent：A→B,  ¬A  ⟹  ¬B（无效）
	•	Bad-DS（错误析取推理）：A∨B,  A  ⟹  ¬B 或 A∨B,  B ⟹ ¬A
	•	Conjunction Intro 误用：A  ⟹  A∧B（缺少对 B 的支持）
	•	Premise Missing：A,  B→C  ⟹  C（漏掉 B）
	•	Direction Flip：A→B  ⟹  B→A（把方向搞反）

一律用语义校验兜底：上面每条都要核验“在某些前提全真情形下结论为假”，避免偶然变真（比如前提不相容）。

⸻

2）生成流程（Generator 内部）

2.1 采样“链长”与模板形态
	•	从 ChainSteps(target) 抽链长 L（如 Easy 多为 1，Hard 可能 2–3）。
	•	在对应的模板集合中等概率抽一个骨架（上面列的 1/2/3 步示例）。
	•	变量分配：从 ["p","q","r","s","t"] 取所需数量的不同变量；Vars 不够时允许复用但尽量避免（Hard 可用更多变量）。

2.2 构造“有效”实例
	•	按骨架构造 Premises 与“主结论”C*。
	•	非空前提域校验：确保存在赋值使所有前提为真（否则换变量或换模板）。
	•	有效性核验：Derivable(Premises, C*) == true。失败就重生。

2.3 扩充 ValidConclusions（凑数量）

在不引入 selection policy 的前提下，用以下保真扩展填充（都要语义核验）：
	•	等价包装：C*、¬¬C*、C* ∨ X（Addition），X 可取未用变量或其否定；
	•	结合已有前提（若有）：A ∧ C*（Conjunction Introduction，前提中含 A）；
	•	链式中间结点：在 A→B, B→C, A ⟹ C 中，也可把 B 作为可推出结论（加核验）。

同时过滤：恒真/恒假结论（真值表全 T / 全 F）不要；字符串去重。

目标：凑够 {SC:≥1,  MC:k∈{2,3}} 的有效结论数量。

2.4 生成 InvalidConclusions（凑数量）
	•	优先用谬误模板按当前 Premises 构造结论；逐条核验 Derivable == false（且存在前提全真赋值）。
	•	不足时，用随机轻微扰动（改联结词、方向翻转、错加/错去一个否定、把必要的前提变量换成未出现的变量等）生成“看似相关”的结论，并核验为不可推出。
	•	同样过滤恒真/恒假、去重。

目标：凑够 {SC:≥3,  MC:(4-k)} 的无效结论数量。

2.5 可行性检查与重试策略（Generator 内）
	•	若 Valid/Invalid 数量达标 → 返回 InferencePools。
	•	否则：
	1.	换骨架重来（同链长）；
	2.	若多次失败 → 放宽链长（例如 L=2 失败则改 L=1 或 3）；
	3.	仍不足 → 换变量实例化或重抽结论扩展；
	4.	达到 MaxRegenAttempts → 上抛 ErrInsufficientCandidates。

这些重试在 Generator 内完成。Orchestrator 收到错误，按你已有“MC=3→2→SC”的回退继续。

⸻

3）与 AllowedOps 的关系
	•	基于模板，不受AllowOps限制

⸻

4）与 Intent/Formatter 的对接
	•	INF_DERIVABLE：正确池=ValidConclusions，干扰池=InvalidConclusions
	•	SC：取 1 + 3；MC：取 k + (4-k)
	•	INF_UNDERIVABLE：池对调
	•	INF_VALIDITY_TF：给定一个结论 C?（可从 Valid/Invalid 里均匀抽），二选一 True/False
	•	dataForTemplate：
	•	SC/MC：{"Premises": numbered_string, ...}（选项里放每个结论字符串）
	•	TF：{"Premises": ..., "C": conclusion_str}

⸻

5）Blueprint 建议
	•	premises_str[]、valid[]、invalid[]
	•	skeleton: 骨架名（MP/MT/DS/HS/Chain2/Chain3…）
	•	chain_steps: L
	•	vars[]: 使用的变量顺序
	•	filters: { removed_tautologies: n, removed_duplicates: n }
	•	（可选）无效结论的 fallacy_tag（AC/DA/BadDS/ConjIntro/DirFlip/MissingPremise）

⸻

6）单测清单（Inference）
	1.	有效性核验
	•	对每个骨架，多次随机变量实例化都能使 Derivable(P,C*)==true，且“前提全真”情形非空。
	2.	谬误核验
	•	对每个 fallacy，随机实例化后 Derivable==false，且存在“前提全真且结论为假”的赋值。
	3.	池可行性
	•	在 Easy/Medium/Hard 下，连跑 1000 次：SC 达到 “1 valid + 3 invalid”；MC(k) 达到 “k valid + (4-k) invalid”，失败率应极低。
	4.	外形与难度一致
	•	在 Easy：输出字符串不包含 →（都被渲染成 ¬A∨B）；
	•	在 Hard：允许出现 →。
	5.	端到端
	•	固定 seed，INF_DERIVABLE + SC 可复现同一题；MC 在 {2,3} 取值不同导致正确项数量不同但均可行。

⸻

7）实现建议（落地顺序）
	1.	先把 Derivable 在 Validator 里写好（真值表 + 非空前提域判断）。
	2.	写一个骨架实例化器：输入骨架与变量，吐出 Premises 与主结论 C*。
	3.	写“有效扩展器”（等价包装、Addition、与前提合取）与“无效生成器”（谬误 + 轻微扰动），每次都调用 Derivable 过滤。
	4.	组装 Generator：循环凑数量 → 返回 pools。
	5.	与 Formatter/Intent 接口对接（装配选项、题干模板填充）。