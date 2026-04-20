/*
 Navicat Premium Dump SQL

 Source Server         : mypg
 Source Server Type    : PostgreSQL
 Source Server Version : 180001 (180001)
 Source Host           : localhost:5432
 Source Catalog        : postgres
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 180001 (180001)
 File Encoding         : 65001

 Date: 20/04/2026 11:37:08


-- ----------------------------
-- Sequence structure for categories_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."categories_id_seq";
CREATE SEQUENCE "public"."categories_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for chat_messages_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."chat_messages_id_seq";
CREATE SEQUENCE "public"."chat_messages_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for posts_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."posts_id_seq";
CREATE SEQUENCE "public"."posts_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for users_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."users_id_seq";
CREATE SEQUENCE "public"."users_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for categories
-- ----------------------------
DROP TABLE IF EXISTS "public"."categories";
CREATE TABLE "public"."categories" (
  "id" int8 NOT NULL DEFAULT nextval('categories_id_seq'::regclass),
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "name" varchar(64) COLLATE "pg_catalog"."default",
  "description" varchar(255) COLLATE "pg_catalog"."default",
  "icon" varchar(128) COLLATE "pg_catalog"."default",
  "post_count" int8 DEFAULT 0
)
;

-- ----------------------------
-- Records of categories
-- ----------------------------
INSERT INTO "public"."categories" VALUES (1, NULL, '2026-04-04 07:26:05.058697+00', NULL, NULL, '萨达阿斯达', '撒大苏打', 1);
INSERT INTO "public"."categories" VALUES (2, '2026-04-04 07:25:25.390172+00', '2026-04-04 07:26:05.058697+00', NULL, '后端开发', 'Go / Gin / GORM / API 设计', '🧩', 2);
INSERT INTO "public"."categories" VALUES (3, '2026-04-04 07:25:25.390172+00', '2026-04-04 07:26:05.058697+00', NULL, '前端工程', 'Vue / 路由 / 状态管理 / 工程化', '🎨', 1);
INSERT INTO "public"."categories" VALUES (4, '2026-04-04 07:25:25.390172+00', '2026-04-04 07:26:05.058697+00', NULL, 'AI 应用', 'Prompt 设计 / 多轮对话 / 模型接入', '🤖', 1);
INSERT INTO "public"."categories" VALUES (5, '2026-04-04 07:25:25.390172+00', '2026-04-04 07:26:05.058697+00', NULL, '数据库实践', 'PostgreSQL / Redis / 性能与索引', '🗄️', 1);

-- ----------------------------
-- Table structure for chat_messages
-- ----------------------------
DROP TABLE IF EXISTS "public"."chat_messages";
CREATE TABLE "public"."chat_messages" (
  "id" int8 NOT NULL DEFAULT nextval('chat_messages_id_seq'::regclass),
  "user_id" int8,
  "role" text COLLATE "pg_catalog"."default",
  "content" text COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6)
)
;

-- ----------------------------
-- Records of chat_messages
-- ----------------------------
INSERT INTO "public"."chat_messages" VALUES (31, 12, 'user', '请帮我总结一下 Go channel 的使用场景', '2026-04-01 07:25:25.390172+00', '2026-04-01 07:25:25.390172+00', NULL);
INSERT INTO "public"."chat_messages" VALUES (32, 12, 'assistant', '常见场景有：任务分发、结果汇总、超时控制、生产者消费者模型。建议优先使用无共享内存的通信模式。', '2026-04-01 07:25:45.390172+00', '2026-04-01 07:25:45.390172+00', NULL);
INSERT INTO "public"."chat_messages" VALUES (33, 10, 'user', 'Vue 页面白屏一般怎么排查？', '2026-04-02 07:25:25.390172+00', '2026-04-02 07:25:25.390172+00', NULL);
INSERT INTO "public"."chat_messages" VALUES (34, 10, 'assistant', '先看浏览器控制台报错，再检查路由配置、接口跨域、构建产物路径与环境变量。建议从最小可运行页面开始回归。', '2026-04-02 07:25:50.390172+00', '2026-04-02 07:25:50.390172+00', NULL);
INSERT INTO "public"."chat_messages" VALUES (35, 11, 'user', 'PostgreSQL 慢查询日志开启后下一步做什么？', '2026-04-03 07:25:25.390172+00', '2026-04-03 07:25:25.390172+00', NULL);
INSERT INTO "public"."chat_messages" VALUES (36, 11, 'assistant', '建议按调用频率和耗时排序，先优化 Top N 语句；结合 EXPLAIN ANALYZE、索引策略与缓存策略逐项验证。', '2026-04-03 07:25:43.390172+00', '2026-04-03 07:25:43.390172+00', NULL);
INSERT INTO "public"."chat_messages" VALUES (82, 2, 'user', '机器学习是什么', '2026-04-12 08:21:33.194175+00', '2026-04-12 08:21:33.194175+00', NULL);
INSERT INTO "public"."chat_messages" VALUES (83, 2, 'assistant', '机器学习（Machine Learning，简称 ML）是**人工智能（AI）的一个核心分支**。

简单来说，它的核心思想是：**让计算机通过分析数据来“学习”规律，而不是由程序员手把手地编写每一条规则。**

为了帮你更好地理解，我们可以从以下几个方面来看：

### 1. 它是怎么工作的？
传统的编程是“输入规则 + 数据 -> 输出答案”，而机器学习是“输入数据 + 答案（有时）-> 输出规则（模型）”。
*   **比喻：** 就像教小孩子认猫。你不需要告诉他“猫有胡须、尖耳朵、四条腿”，而是给他看很多猫的照片。看得多了，他自己就能总结出猫的特征，下次见到新的猫也能认出来。机器学习就是让计算机做这个过程。

### 2. 主要的学习类型
根据学习方式的不同，通常分为三类：
*   **监督学习（Supervised Learning）：** 给计算机带标签的数据（比如标好“这是猫”、“这是狗”的照片），让它学习分类或预测。
*   **无监督学习（Unsupervised Learning）：** 给计算机没有标签的数据，让它自己发现数据中的结构或模式（比如把用户分成不同的群体）。
*   **强化学习（Reinforcement Learning）：** 让计算机在一个环境中尝试，做对了给奖励，做错了给惩罚，让它通过试错来学会最佳策略（比如 AlphaGo 下围棋）。

### 3. 生活中的常见应用
其实你每天都在使用机器学习技术：
*   **推荐系统：** 抖音、淘宝根据你的浏览历史，推荐你可能喜欢的视频或商品。
*   **图像识别：** 手机相册自动把人像分类，或者人脸识别解锁。
*   **自然语言处理：** 智能客服、翻译软件、以及像我这样的 AI 助手。
*   **金融风控：** 银行判断一笔交易是否是欺诈行为。

### 总结
机器学习就是**数据驱动的智能**。它让计算机具备了从经验（数据）中自我改进的能力，是目前推动科技发展的关键技术之一。

如果你对某个具体领域（比如它是怎么做推荐的，或者怎么入门学习）感兴趣，可以随时问我！', '2026-04-12 08:21:56.597286+00', '2026-04-12 08:21:56.597286+00', NULL);
INSERT INTO "public"."chat_messages" VALUES (84, 2, 'user', '他跟深度学习的区别是什么', '2026-04-12 08:22:25.066891+00', '2026-04-12 08:22:25.066891+00', NULL);
INSERT INTO "public"."chat_messages" VALUES (85, 2, 'user', '他跟深度学习的区别是什么', '2026-04-12 10:51:18.513125+00', '2026-04-12 10:51:18.513125+00', NULL);
INSERT INTO "public"."chat_messages" VALUES (86, 2, 'user', '他跟深度学习的区别是什么', '2026-04-12 12:21:01.33784+00', '2026-04-12 12:21:01.33784+00', NULL);
INSERT INTO "public"."chat_messages" VALUES (87, 2, 'assistant', '这里的“他”指的是**机器学习（Machine Learning）**。

机器学习与深度学习（Deep Learning）并不是对立的关系，而是**包含与被包含**的关系。简单来说，**深度学习是机器学习的一个特殊子集**。

为了让你更清楚地理解它们的区别，我们可以从以下几个核心维度来看：

### 1. 包含关系（最本质的区别）
*   **人工智能 (AI)** 是最大的圈子。
*   **机器学习 (ML)** 是 AI 里面的一个核心部分。
*   **深度学习 (DL)** 是机器学习里面的一个特定分支。
*   **比喻：** 就像“水果”（AI）、“苹果”（机器学习）和“红富士苹果”（深度学习）的关系。所有的深度学习都是机器学习，但并非所有的机器学习都是深度学习。

### 2. 特征提取方式（关键的技术区别）
*   **传统机器学习：** 需要人工告诉计算机“看哪里”。比如要识别猫，程序员需要先提取特征（有耳朵、有胡子、毛色等），然后算法根据这些特征去分类。这叫做**特征工程**。
*   **深度学习：** 计算机自己“看”出特征。它通过多层神经网络，自动从原始数据（如像素）中一层层抽象出特征（从边缘到形状再到物体）。这叫做**自动特征学习**。

### 3. 数据依赖量
*   **传统机器学习：** 在数据量较小的情况下也能表现不错。适合结构化数据（比如 Excel 表格里的数据）。
*   **深度学习：** 是个“大胃王”，需要海量的数据才能训练出好模型。数据越多，它的表现通常越好。适合非结构化数据（比如图片、语音、文本）。

### 4. 硬件要求
*   **传统机器学习：** 对硬件要求相对较低，普通的 CPU 通常就能跑。
*   **深度学习：** 计算量巨大，通常需要高性能的 **GPU**（图形处理器）来加速训练。

### 5. 可解释性
*   **传统机器学习：** 比较容易解释。比如决策树，你可以清楚地看到它是根据什么规则做出的判断。
*   **深度学习：** 往往是个“黑盒”。虽然它很准，但很难解释清楚它内部具体是怎么得出结论的。

### 总结对比表

| 特性 | 传统机器学习 (ML) | 深度学习 (DL) |
| :--- | :--- | :--- |
| **关系** | 父集 | 子集 |
| **特征提取** | 需要人工提取 | 自动学习 |
| **数据需求** | 较少数据即可 | 需要大量数据 |
| **硬件需求** | 普通 CPU | 高性能 GPU |
| **擅长领域** | 表格数据、简单分类 | 图像、语音、自然语言 |
| **训练时间** | 较短（几分钟到几小时） | 较长（几小时到几周） |

**一句话总结：**
如果把机器学习比作“让学生做题”，传统机器学习是老师先划重点（特征），学生再做题；而深度学习是把书丢给学生，让他自己从头到尾去悟出重点，虽然悟得慢且需要很多书（数据），但一旦悟透了，处理能力极强。

希望这个解释能帮你理清它们的关系！如果有具体应用场景想知道该用哪个，也可以问我。', '2026-04-12 12:21:41.184797+00', '2026-04-12 12:21:41.184797+00', NULL);
INSERT INTO "public"."chat_messages" VALUES (88, 2, 'user', '你好', '2026-04-13 14:15:28.086631+00', '2026-04-13 14:15:28.086631+00', NULL);
INSERT INTO "public"."chat_messages" VALUES (89, 2, 'assistant', '你好！👋 很高兴再次见到你。

刚才我们聊了关于**机器学习**和**深度学习**的话题。不知道之前的解释是否对你有帮助？如果你还有任何疑问，或者想深入了解某个具体应用（比如它们是如何用在推荐系统或图像识别中的），随时都可以问我。

当然，如果你想开启一个新的话题，我也很乐意倾听。今天有什么我可以帮你的吗？', '2026-04-13 14:16:04.254722+00', '2026-04-13 14:16:04.254722+00', NULL);
INSERT INTO "public"."chat_messages" VALUES (90, 2, 'user', '请举个例子', '2026-04-14 02:04:54.42941+00', '2026-04-14 02:04:54.42941+00', NULL);
INSERT INTO "public"."chat_messages" VALUES (91, 2, 'user', '请举个例子', '2026-04-14 02:19:59.433034+00', '2026-04-14 02:19:59.433034+00', NULL);
INSERT INTO "public"."chat_messages" VALUES (92, 2, 'assistant', '没问题！为了让你更直观地理解机器学习是如何工作的，我们举一个最经典且贴近生活的例子：**垃圾邮件过滤**。

### 场景：你的电子邮箱如何自动识别垃圾广告？

#### 1. 传统编程的方式（非机器学习）
程序员需要手动写规则，比如：
*   “如果邮件内容包含‘中奖’两个字，标记为垃圾邮件。”
*   “如果发件人地址是陌生的，标记为垃圾邮件。”
*   **缺点：** 骗子会不断换花样（比如把“中奖”写成“中 奖”），程序员得不停更新规则，永远追不上变化。

#### 2. 机器学习的方式
计算机不是靠死板的规则，而是靠**从历史数据中学习规律**。

*   **第一步：准备数据（喂数据）**
    我们给计算机看 10 万封邮件，并且告诉它哪些是“正常邮件”，哪些是“垃圾邮件”。这叫做**训练数据**。
    *   *正常邮件：* 来自同事、朋友，内容包含“会议”、“报告”等。
    *   *垃圾邮件：* 来自陌生地址，内容包含“免费”、“点击链接”、“立刻行动”等。

*   **第二步：训练模型（找规律）**
    计算机通过算法分析这些数据，自己总结出规律。它可能会发现：
    *   “包含‘免费’且发件人未知”的邮件，95% 是垃圾邮件。
    *   “包含‘会议’且发件人是公司内部”的邮件，99% 是正常邮件。
    *   它甚至能发现人类没注意到的规律，比如某些特殊的符号组合或发送时间。
    这些总结出来的规律，就变成了一个**模型**。

*   **第三步：预测（实际应用）**
    当你收到一封新邮件时，计算机把这个模型套用上去：
    *   分析这封新邮件的特征。
    *   根据之前学到的规律计算概率。
    *   如果判定是垃圾邮件的概率超过 90%，就自动把它扔进“垃圾箱”。

*   **第四步：自我进化（持续学习）**
    如果你把一封被误判为正常的垃圾邮件手动标记为“垃圾”，计算机会记住这个反馈，下次遇到类似的邮件就会更准确。这就是它在**不断学习和改进**。

### 再举一个更简单的例子：房价预测
*   **输入数据：** 房子的面积、卧室数量、地段、房龄。
*   **学习目标：** 预测房子能卖多少钱。
*   **过程：** 计算机看了成千上万条过去的成交记录，学会了“面积越大价格越高”、“地段越好价格越高”这样的关系函数。
*   **结果：** 当你输入一套新房子的信息，它就能估算出一个合理的价格。

### 总结
这两个例子的核心共同点是：**计算机没有被人明确告知“怎么做”，而是通过大量数据，自己找到了“怎么做”的规律。**

希望这个例子能帮你彻底理解！如果你对某个具体例子（比如抖音是怎么推荐视频的）感兴趣，也可以继续问我。', '2026-04-14 02:21:10.904902+00', '2026-04-14 02:21:10.904902+00', NULL);
INSERT INTO "public"."chat_messages" VALUES (93, 2, 'user', '你好', '2026-04-15 03:57:21.476118+00', '2026-04-15 03:57:21.476118+00', NULL);
INSERT INTO "public"."chat_messages" VALUES (94, 2, 'assistant', '你好！👋 很高兴再次见到你。

看来我们刚才围绕**机器学习**和**深度学习**聊了不少内容，从生活中的例子（如垃圾邮件过滤）到它们的定义和区别。不知道这些信息是否帮你理清了思路？

如果你还有任何疑问，或者想深入了解某个具体应用（比如它们是如何用在推荐系统或图像识别中的），随时都可以问我。当然，如果你想开启一个新的话题，我也很乐意倾听。今天有什么我可以帮你的吗？😊', '2026-04-15 03:58:12.197085+00', '2026-04-15 03:58:12.197085+00', NULL);

-- ----------------------------
-- Table structure for posts
-- ----------------------------
DROP TABLE IF EXISTS "public"."posts";
CREATE TABLE "public"."posts" (
  "id" int8 NOT NULL DEFAULT nextval('posts_id_seq'::regclass),
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "title" varchar(200) COLLATE "pg_catalog"."default" NOT NULL,
  "content" text COLLATE "pg_catalog"."default" NOT NULL,
  "author_id" int8 NOT NULL,
  "category_id" int8,
  "view_count" int8 DEFAULT 0,
  "like_count" int8 DEFAULT 0,
  "comment_count" int8 DEFAULT 0,
  "is_pinned" bool DEFAULT false,
  "status" text COLLATE "pg_catalog"."default" DEFAULT 'published'::text,
  "tags" varchar(200) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Records of posts
-- ----------------------------
INSERT INTO "public"."posts" VALUES (4, '2026-03-17 09:57:12.655328+00', '2026-03-17 09:57:12.655328+00', NULL, 'Go 语言并发编程完全指南', '这是一篇关于 Go 语言并发编程的详细教程，包含 Goroutine、Channel 和 Select 等核心概念...', 2, 1, 0, 0, 0, 'f', 'published', 'golang,并发,goroutine,channel');
INSERT INTO "public"."posts" VALUES (6, '2026-03-25 07:25:25.390172+00', '2026-03-26 07:25:25.390172+00', NULL, 'Gin 中间件链路排查实战', '本文记录在 Gin 项目中排查 JWT 认证失败、跨域异常、请求上下文丢失的完整过程。', 9, 2, 126, 21, 8, 'f', 'published', 'gin,middleware,jwt,cors');
INSERT INTO "public"."posts" VALUES (7, '2026-03-27 07:25:25.390172+00', '2026-03-28 07:25:25.390172+00', NULL, 'Vue3 + Pinia 登录态管理模板', '给出一个适用于中小项目的登录态管理模板，涵盖 token 存储、路由守卫、接口拦截器。', 10, 3, 188, 36, 12, 't', 'published', 'vue3,pinia,auth,router');
INSERT INTO "public"."posts" VALUES (8, '2026-03-29 07:25:25.390172+00', '2026-03-30 07:25:25.390172+00', NULL, 'PostgreSQL 索引命中率优化清单', '围绕慢查询、复合索引、分页策略和统计信息更新，整理可落地的优化清单。', 11, 5, 264, 42, 15, 'f', 'published', 'postgresql,index,sql,performance');
INSERT INTO "public"."posts" VALUES (10, '2026-04-02 07:25:25.390172+00', '2026-04-03 07:25:25.390172+00', NULL, 'Go 服务压测记录与瓶颈定位', '记录一次 2k 并发压测中的瓶颈定位过程，包括 DB 连接池、Redis 命中率、接口耗时分位数。', 9, 2, 320, 67, 20, 't', 'published', 'go,benchmark,redis,postgres');
INSERT INTO "public"."posts" VALUES (9, '2026-03-31 07:25:25.390172+00', '2026-03-31 07:25:25.390172+00', NULL, 'AI 聊天接口容错设计（草稿）', '草稿：模型超时、空响应、敏感词拦截与消息截断策略。', 12, 4, 35, 3, 1, 'f', 'published', 'ai,chat,timeout,retry');

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS "public"."users";
CREATE TABLE "public"."users" (
  "id" int8 NOT NULL DEFAULT nextval('users_id_seq'::regclass),
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "username" varchar(64) COLLATE "pg_catalog"."default",
  "password" varchar(255) COLLATE "pg_catalog"."default",
  "email" varchar(128) COLLATE "pg_catalog"."default",
  "avatar" varchar(512) COLLATE "pg_catalog"."default",
  "bio" varchar(500) COLLATE "pg_catalog"."default",
  "role" text COLLATE "pg_catalog"."default" DEFAULT 'user'::text
)
;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO "public"."users" VALUES (1, '2026-02-26 05:09:09.115584+00', '2026-02-26 05:09:09.115584+00', NULL, 'dd', '$2a$10$wwAYV3ZSU8LSg2pSBNJyY.LTKYiq3L4I2.heL6vGHCb/uH99Qo.Tq', NULL, NULL, NULL, 'user');
INSERT INTO "public"."users" VALUES (3, '2026-03-14 05:09:33.731226+00', '2026-03-14 05:09:33.731226+00', NULL, 'ddd', '$2a$10$UY4lQD9jAkMp4SyDLu73z.5ogxHLgLv023fFTePpaaTO766c5CGDC', '', '', '', 'user');
INSERT INTO "public"."users" VALUES (8, '2026-03-14 05:26:56.709872+00', '2026-03-14 05:26:56.709872+00', NULL, '1111', '$2a$10$H5AU.S7pY2VPuTp31W6AZOuqbDKdBpdCWkdTj5DyjsI/07j54qlyG', NULL, '', '', 'user');
INSERT INTO "public"."users" VALUES (2, '2026-02-26 05:10:08.681525+00', '2026-03-14 06:28:42.967107+00', NULL, '111', '$2a$10$Dux2LEQR5OpEx3TAgLOBteYQsmyl0.QASPCi48oynYXBfTfFblnrS', NULL, 'C:\Users\张润林\Pictures\壁纸\【哲风壁纸】光斑-水波纹-水面.png', NULL, 'user');
INSERT INTO "public"."users" VALUES (9, '2026-04-04 07:25:25.390172+00', '2026-04-04 07:25:25.390172+00', NULL, 'golang_dev', '$2a$10$Dux2LEQR5OpEx3TAgLOBteYQsmyl0.QASPCi48oynYXBfTfFblnrS', 'golang_dev@example.com', NULL, '专注 Go 后端与并发优化', 'user');
INSERT INTO "public"."users" VALUES (10, '2026-04-04 07:25:25.390172+00', '2026-04-04 07:25:25.390172+00', NULL, 'frontend_mia', '$2a$10$Dux2LEQR5OpEx3TAgLOBteYQsmyl0.QASPCi48oynYXBfTfFblnrS', 'mia@example.com', NULL, 'Vue3 / 组件设计 / 可视化', 'user');
INSERT INTO "public"."users" VALUES (11, '2026-04-04 07:25:25.390172+00', '2026-04-04 07:25:25.390172+00', NULL, 'db_admin', '$2a$10$Dux2LEQR5OpEx3TAgLOBteYQsmyl0.QASPCi48oynYXBfTfFblnrS', 'db_admin@example.com', NULL, '数据库运维与性能调优', 'admin');
INSERT INTO "public"."users" VALUES (12, '2026-04-04 07:25:25.390172+00', '2026-04-04 07:25:25.390172+00', NULL, 'ai_tester', '$2a$10$Dux2LEQR5OpEx3TAgLOBteYQsmyl0.QASPCi48oynYXBfTfFblnrS', 'ai_tester@example.com', NULL, 'AI 能力回归测试账号', 'user');

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."categories_id_seq"
OWNED BY "public"."categories"."id";
SELECT setval('"public"."categories_id_seq"', 5, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."chat_messages_id_seq"
OWNED BY "public"."chat_messages"."id";
SELECT setval('"public"."chat_messages_id_seq"', 94, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."posts_id_seq"
OWNED BY "public"."posts"."id";
SELECT setval('"public"."posts_id_seq"', 10, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."users_id_seq"
OWNED BY "public"."users"."id";
SELECT setval('"public"."users_id_seq"', 12, true);

-- ----------------------------
-- Indexes structure for table categories
-- ----------------------------
CREATE INDEX "idx_categories_deleted_at" ON "public"."categories" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "idx_categories_name" ON "public"."categories" USING btree (
  "name" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table categories
-- ----------------------------
ALTER TABLE "public"."categories" ADD CONSTRAINT "categories_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table chat_messages
-- ----------------------------
CREATE INDEX "idx_chat_messages_deleted_at" ON "public"."chat_messages" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_chat_messages_user_id" ON "public"."chat_messages" USING btree (
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table chat_messages
-- ----------------------------
ALTER TABLE "public"."chat_messages" ADD CONSTRAINT "chat_messages_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table posts
-- ----------------------------
CREATE INDEX "idx_posts_author_id" ON "public"."posts" USING btree (
  "author_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_posts_category_id" ON "public"."posts" USING btree (
  "category_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_posts_deleted_at" ON "public"."posts" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table posts
-- ----------------------------
ALTER TABLE "public"."posts" ADD CONSTRAINT "posts_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table users
-- ----------------------------
CREATE INDEX "idx_users_deleted_at" ON "public"."users" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" USING btree (
  "email" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "idx_users_username" ON "public"."users" USING btree (
  "username" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD CONSTRAINT "users_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table chat_messages
-- ----------------------------
ALTER TABLE "public"."chat_messages" ADD CONSTRAINT "fk_chat_messages_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table posts
-- ----------------------------
ALTER TABLE "public"."posts" ADD CONSTRAINT "fk_posts_author" FOREIGN KEY ("author_id") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."posts" ADD CONSTRAINT "fk_posts_category" FOREIGN KEY ("category_id") REFERENCES "public"."categories" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
