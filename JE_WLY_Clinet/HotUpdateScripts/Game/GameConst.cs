using HotUpdateScripts.ProtoM;
using System.Collections.Generic;

using UnityEngine;


namespace HotUpdateScripts
{
    public class GameConst
    {
        public enum ModelAnim
        {
            No,
            attack,
            walk,
            run,
            dead,
            stand,
            skill,
            skill1,
            skill2,
            skill3,
            skill4,
            skill5,
            shouji_bingfeng,
            shouji_chenmo,
            shouji_hunluan,
            shouji_kunbang,
            shouji_xuanyun,
            shouji_zhongdu,
        }

        private static Dictionary<int, ModelAnim> SkillNeedToAmim = new Dictionary<int, ModelAnim>
        {
            { 0, ModelAnim.attack },
            { 1, ModelAnim.skill1 },
            { 2, ModelAnim.skill2 },
            { 3, ModelAnim.skill3 },
            { 4, ModelAnim.skill4 },
            { 5, ModelAnim.skill5 },
            { 100, ModelAnim.skill },
        };

        public static ModelAnim GetAnim(int idx)
        {
            if (SkillNeedToAmim.TryGetValue(idx, out var ani))
                return ani;

            return ModelAnim.No;
        }

        // 自动开箱选项
        public enum AutoOpenBoxOption
        {
            Quality,        // 大于等于这个品质
            QualityToggle,  // 是否根据品质筛选
            Level,          // 大于等于这个等级
            LevelToggle,    // 是否根据等级筛选
            Occupation,     // 筛选的职业
            OccupationToggle,   // 是否根据职业筛选
            BagFullToggle,  //  背包是否已满
            HigherThenProtagonistEquippedToggle,    // 获得比主角当前穿戴装备战力更高的装备
            HigherThenPartnerEquippedToggle,        // 获得比伙伴当前穿戴装备战力更高的装备
            DefaultOpenSoltCondition, // 默认开出哪个槽位的装备
            AutoRecycleToggle,  // 不满足上述条件的装备自动回收
        }

        public enum ChatOption
        {
            WorldVoice,
            GuildVoice,
            WorldMSG,
            GuildMSG,
        }

        public enum SystemSettingOption
        {
            BgMusic,
            SoundEffect,
            DisplayMode,    // 显示模式
            ImageQuality,   // 画质
            TargetFrameRate,// 帧率
            EffectSetting,  // 特效显示
        }

        public static int GetNormallSkillByJob(JobType jobType)
        {
            switch (jobType)
            {
                case JobType.Warrior:
                    return 100;
                case JobType.Mage:
                    return 1400;
                case JobType.Sacrifice:
                    return 2700;
                default:
                    return 100;
            }
        }

        private static Dictionary<string, ModelClass> MainRolePathDic = new Dictionary<string, ModelClass>
        {
            [$"{JobType.Warrior}{SexType.SexGril}"] = new ModelClass { Path = "Role/Prefabs/Warrior_Gril", BornPatch = "Other/Effects/GameSpine/Zs_bron_g/Bron" },
            [$"{JobType.Warrior}{SexType.SexBoy}"] = new ModelClass { Path = "Role/Prefabs/Warrior_Boy", BornPatch = "Other/Effects/GameSpine/Zs_bron/Warrior_Bron" },
            [$"{JobType.Warrior}{SexType.SexGril}_Fight"] = new ModelClass { Path = "Role/Prefabs/Fight_Warrior_Gril" },
            [$"{JobType.Warrior}{SexType.SexBoy}_Fight"] = new ModelClass { Path = "Role/Prefabs/Fight_Warrior_Boy" },

            [$"{JobType.Mage}{SexType.SexGril}"] = new ModelClass { Path = "Role/Prefabs/Mage_Gril", BornPatch = "Other/Effects/GameSpine/Fs_bron_g/Bron" },
            [$"{JobType.Mage}{SexType.SexBoy}"] = new ModelClass { Path = "Role/Prefabs/Mage_Boy", BornPatch = "Other/Effects/GameSpine/Fs_bron/Mage_Bron" },
            [$"{JobType.Mage}{SexType.SexGril}_Fight"] = new ModelClass { Path = "Role/Prefabs/Fight_Mage_Gril" },
            [$"{JobType.Mage}{SexType.SexBoy}_Fight"] = new ModelClass { Path = "Role/Prefabs/Fight_Mage_Boy" },
            [$"{JobType.Mage}{SexType.SexBoy}_Fight"] = new ModelClass { Path = "Role/Prefabs/Fight_Mage_Boy" },

            [$"{JobType.Sacrifice}{SexType.SexGril}"] = new ModelClass { Path = "Role/Prefabs/Sacrifice_Gril", BornPatch = "Other/Effects/GameSpine/Js_bron_g/Bron" },
            [$"{JobType.Sacrifice}{SexType.SexBoy}"] = new ModelClass { Path = "Role/Prefabs/Sacrifice_Boy", BornPatch = "Other/Effects/GameSpine/Js_bron/Sacrifice_Bron" },
            [$"{JobType.Sacrifice}{SexType.SexGril}_Fight"] = new ModelClass { Path = "Role/Prefabs/Fight_Sacrifice_Gril" },
            [$"{JobType.Sacrifice}{SexType.SexBoy}_Fight"] = new ModelClass { Path = "Role/Prefabs/Fight_Sacrifice_Boy" },
        };

        public static string GetModelByJobAndSex(JobType jobType, SexType sexType, bool isFight = false)
        {
            string poolKey = $"{jobType}{sexType}";
            if (isFight)
                poolKey += "_Fight";
            return BaseResousePath + MainRolePathDic[poolKey].Path + ".prefab";
        }

        public static string GetBornPatchByJobAndSex(JobType jobType, SexType sexType)
        {
            string poolKey = $"{jobType}{sexType}";
            if (MainRolePathDic[poolKey].BornPatch == null)
                return string.Empty;

            return BaseResousePath + MainRolePathDic[poolKey].BornPatch + ".prefab";
        }

        public static string GetModelPath(string path)
        {
            if (!path.Contains(".prefab"))
            {
                path += ".prefab";
            }
            if (!path.Contains(BaseResousePath))
            {
                path = BaseResousePath + path;
            }
            return path;
        }

        private struct ModelClass
        {
            public string Path;
            public float Scal;
            public string BornPatch;
        }

        private static Dictionary<string, ModelClass> MonsterDic = new Dictionary<string, ModelClass>
        {
            ["Boss_1"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Boss_1", Scal = 0.9f },
            ["Boss_2"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Boss_2", Scal = 0.75f },
            ["Boss_3"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Boss_3", Scal = 0.75f },
            ["Boss_4"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Boss_4", Scal = 0.825f },
            ["Boss_5"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Boss_5", Scal = 0.825f },
            ["Boss_6"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Boss_6", Scal = 0.825f },
            ["Boss_7"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Boss_7", Scal = 0.825f },
            ["Boss_8"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Boss_8", Scal = 0.825f },
            ["Boss_9"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Boss_9", Scal = 0.825f },
            ["Boss_9_Call"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Boss_9_Call", Scal = 0.825f },
            ["Boss_Romate_1"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Boss_Romate_1", Scal = 0.825f },
            ["Boss_Romate_2"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Boss_Romate_2", Scal = 0.9f },
            ["Boss_Romate_3"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Boss_Romate_3", Scal = 0.9f },
            ["Boss_Romate_4"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Boss_Romate_4", Scal = 0.825f },
            ["Boss_Romate_5"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Boss_Romate_5", Scal = 0.9f },
            ["Boss_Romate_6"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Boss_Romate_6", Scal = 0.9f },
            ["Boss_Romate_7"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Boss_Romate_7", Scal = 0.9f },

            //小怪
            ["Monster_1"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Monster_1", Scal = 0.495f },
            ["Monster_2"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Monster_2", Scal = 0.495f },
            ["Monster_3"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Monster_3", Scal = 0.495f },
            ["Monster_4"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Monster_4", Scal = 0.495f },
            ["Monster_5"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Monster_5", Scal = 0.495f },
            ["Monster_6"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Monster_6", Scal = 0.495f },
            ["Monster_7"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Monster_7", Scal = 0.495f },
            ["Monster_8"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Monster_8", Scal = 0.495f },
            ["Monster_9"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Monster_9", Scal = 0.45f },
            ["Monster_10"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Monster_10", Scal = 0.495f },
            ["Monster_11"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Monster_11", Scal = 0.495f },
            ["Monster_12"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Monster_12", Scal = 0.495f },
            ["Monster_13"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Monster_13", Scal = 0.495f },
            ["Monster_14"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Monster_14", Scal = 0.495f },
            ["Monster_15"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Monster_15", Scal = 0.495f },
            ["Monster_16"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Monster_16", Scal = 0.495f },

            //远程小怪
            ["Monster_Romate_1"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Ranged _Monster_1", Scal = 0.495f },
            ["Monster_Romate_2"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Ranged _Monster_2", Scal = 0.495f },
            ["Monster_Romate_3"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Ranged _Monster_3", Scal = 0.495f },
            ["Monster_Romate_4"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Ranged _Monster_4", Scal = 0.495f },
            ["Monster_Romate_5"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Ranged _Monster_5", Scal = 0.495f },
            ["Monster_Romate_6"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Ranged _Monster_6", Scal = 0.495f },
            ["Monster_Romate_7"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Ranged _Monster_7", Scal = 0.495f },
            ["Monster_Romate_8"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Ranged _Monster_8", Scal = 0.495f },
            ["Monster_Romate_9"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Ranged _Monster_9", Scal = 0.495f },
            ["Monster_Romate_10"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Ranged _Monster_10", Scal = 0.495f },

            ["Monster_Box1_Blue"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Box", Scal = 2f },
            ["Monster_Box1_Purple"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Box", Scal = 2f },
            ["Monster_Box1_Yellow"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Box", Scal = 2f },
            ["Monster_Box_Blue"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Box", Scal = 2f },
            ["Monster_Box_Purple"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Box", Scal = 2f },
            ["Monster_Box_Yellow"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_Box", Scal = 2f },

            ["Monster_FuDai"] = new ModelClass { Path = "Role/MonsterPrefabs/Fight_FuDai", Scal = 2f },

            //战斗伙伴
            ["Fight_Partner_1"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_1", Scal = 0.8f },
            ["Fight_Partner_2"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_2", Scal = 0.6f },
            ["Fight_Partner_3"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_3", Scal = 0.65f },
            ["Fight_Partner_4"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_4", Scal = 0.8f },
            ["Fight_Partner_5"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_5", Scal = 0.8f },
            ["Fight_Partner_6"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_6", Scal = 0.75f },
            ["Fight_Partner_7"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_7", Scal = 0.75f },
            ["Fight_Partner_8"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_8", Scal = 0.75f },
            ["Fight_Partner_9"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_9", Scal = 0.75f },
            ["Fight_Partner_10"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_10", Scal = 0.7f },
            ["Fight_Partner_11"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_11", Scal = 0.7f },
            ["Fight_Partner_12"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_12", Scal = 0.76f },
            ["Fight_Partner_13"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_13", Scal = 0.75f },
            ["Fight_Partner_14"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_14", Scal = 0.74f },
            ["Fight_Partner_15"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_15", Scal = 0.75f },
            ["Fight_Partner_16"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_16", Scal = 0.75f },
            ["Fight_Partner_17"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_17", Scal = 0.75f },
            ["Fight_Partner_18"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_18", Scal = 0.7f },
            ["Fight_Partner_19"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_19", Scal = 0.6f },
            ["Fight_Partner_20"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_20", Scal = 0.6f },
            ["Fight_Partner_21"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_21", Scal = 0.6f },
            ["Fight_Partner_22"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_22", Scal = 0.6f },
            ["Fight_Partner_23"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_23", Scal = 0.6f },
            ["Fight_Partner_24"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_24", Scal = 0.6f },
            ["Fight_Partner_25"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_25", Scal = 0.65f },
            ["Fight_Partner_26"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_26", Scal = 0.6f },
            ["Fight_Partner_27"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_27", Scal = 0.65f },
            ["Fight_Partner_28"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_28", Scal = 0.65f },
            ["Fight_Partner_29"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_29", Scal = 0.65f },
            ["Fight_Partner_30"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_30", Scal = 0.65f },
            ["Fight_Partner_31"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_31", Scal = 0.65f },
            ["Fight_Partner_32"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_32", Scal = 0.65f },
            ["Fight_Partner_33"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_33", Scal = 0.7f },
            ["Fight_Partner_51"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_51", Scal = 0.85f },
            ["Fight_Partner_52"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_52", Scal = 0.85f },
            ["Fight_Partner_53"] = new ModelClass { Path = "Role/Prefabs/Fight_Partner_53", Scal = 0.9f },

            //UI伙伴
            ["Partner_1"] = new ModelClass { Path = "Role/Prefabs/Partner_1", Scal = 1 },
            ["Partner_2"] = new ModelClass { Path = "Role/Prefabs/Partner_2", Scal = 1 },
            ["Partner_3"] = new ModelClass { Path = "Role/Prefabs/Partner_3", Scal = 1 },
            ["Partner_4"] = new ModelClass { Path = "Role/Prefabs/Partner_4", Scal = 1 },
            ["Partner_5"] = new ModelClass { Path = "Role/Prefabs/Partner_5", Scal = 1 },
            ["Partner_6"] = new ModelClass { Path = "Role/Prefabs/Partner_6", Scal = 1 },
            ["Partner_7"] = new ModelClass { Path = "Role/Prefabs/Partner_7", Scal = 1 },
            ["Partner_8"] = new ModelClass { Path = "Role/Prefabs/Partner_8", Scal = 1 },
            ["Partner_9"] = new ModelClass { Path = "Role/Prefabs/Partner_9", Scal = 1 },
            ["Partner_10"] = new ModelClass { Path = "Role/Prefabs/Partner_10", Scal = 1 },
            ["Partner_11"] = new ModelClass { Path = "Role/Prefabs/Partner_11", Scal = 1 },
            ["Partner_12"] = new ModelClass { Path = "Role/Prefabs/Partner_12", Scal = 1 },
            ["Partner_13"] = new ModelClass { Path = "Role/Prefabs/Partner_13", Scal = 1 },
            ["Partner_14"] = new ModelClass { Path = "Role/Prefabs/Partner_14", Scal = 1 },
            ["Partner_15"] = new ModelClass { Path = "Role/Prefabs/Partner_15", Scal = 1 },
            ["Partner_16"] = new ModelClass { Path = "Role/Prefabs/Partner_16", Scal = 1 },
            ["Partner_17"] = new ModelClass { Path = "Role/Prefabs/Partner_17", Scal = 1 },
            ["Partner_18"] = new ModelClass { Path = "Role/Prefabs/Partner_18", Scal = 1 },
            ["Partner_19"] = new ModelClass { Path = "Role/Prefabs/Partner_19", Scal = 1 },
            ["Partner_20"] = new ModelClass { Path = "Role/Prefabs/Partner_20", Scal = 1 },
            ["Partner_21"] = new ModelClass { Path = "Role/Prefabs/Partner_21", Scal = 1 },
            ["Partner_22"] = new ModelClass { Path = "Role/Prefabs/Partner_22", Scal = 1 },
            ["Partner_23"] = new ModelClass { Path = "Role/Prefabs/Partner_23", Scal = 1 },
            ["Partner_24"] = new ModelClass { Path = "Role/Prefabs/Partner_24", Scal = 1 },
            ["Partner_25"] = new ModelClass { Path = "Role/Prefabs/Partner_25", Scal = 1 },
            ["Partner_26"] = new ModelClass { Path = "Role/Prefabs/Partner_26", Scal = 1 },
            ["Partner_27"] = new ModelClass { Path = "Role/Prefabs/Partner_27", Scal = 1 },
            ["Partner_28"] = new ModelClass { Path = "Role/Prefabs/Partner_28", Scal = 1 },
            ["Partner_29"] = new ModelClass { Path = "Role/Prefabs/Partner_29", Scal = 1 },
            ["Partner_30"] = new ModelClass { Path = "Role/Prefabs/Partner_30", Scal = 1 },
            ["Partner_31"] = new ModelClass { Path = "Role/Prefabs/Partner_31", Scal = 1 },
            ["Partner_32"] = new ModelClass { Path = "Role/Prefabs/Partner_32", Scal = 1 },
            ["Partner_33"] = new ModelClass { Path = "Role/Prefabs/Partner_33", Scal = 1 },
            ["Partner_51"] = new ModelClass { Path = "Role/Prefabs/Partner_51", Scal = 1 },
            ["Partner_52"] = new ModelClass { Path = "Role/Prefabs/Partner_52", Scal = 1 },
            ["Partner_53"] = new ModelClass { Path = "Role/Prefabs/Partner_53", Scal = 1 },

            //宠物
            ["Pet_1"] = new ModelClass { Path = "Role/Prefabs/Pet_Dog", Scal = 1 },
            ["Pet_2"] = new ModelClass { Path = "Role/Prefabs/Pet_Cat", Scal = 1 },
            ["Pet_3"] = new ModelClass { Path = "Role/Prefabs/Pet_Mouse", Scal = 1 },
            ["Pet_4"] = new ModelClass { Path = "Role/Prefabs/Pet_Llama", Scal = 1 },
            ["Pet_5"] = new ModelClass { Path = "Role/Prefabs/Pet_Rabbit", Scal = 1 },
            ["Pet_6"] = new ModelClass { Path = "Role/Prefabs/Pet_Bear", Scal = 1 },
            ["Pet_7"] = new ModelClass { Path = "Role/Prefabs/Pet_Deer", Scal = 1 },
            ["Pet_8"] = new ModelClass { Path = "Role/Prefabs/Pet_Fox", Scal = 1 },
            ["Pet_9"] = new ModelClass { Path = "Role/Prefabs/Pet_Kangaroo", Scal = 1 },
            ["Pet_10"] = new ModelClass { Path = "Role/Prefabs/Pet_Pig", Scal = 1 },
            ["Pet_11"] = new ModelClass { Path = "Role/Prefabs/Pet_Turtle", Scal = 1 },
            ["Pet_12"] = new ModelClass { Path = "Role/Prefabs/Pet_Exp", Scal = 1 },
            ["Pet_13"] = new ModelClass { Path = "Role/Prefabs/Pet_Exp", Scal = 1 },
            ["Pet_14"] = new ModelClass { Path = "Role/Prefabs/Pet_Exp", Scal = 1 },
            ["Pet_15"] = new ModelClass { Path = "Role/Prefabs/Pet_Exp", Scal = 1 },
            ["Pet_16"] = new ModelClass { Path = "Role/Prefabs/Pet_Exp", Scal = 1 },
            ["Pet_17"] = new ModelClass { Path = "Role/Prefabs/Pet_Exp", Scal = 1 },
            ["Pet_21"] = new ModelClass { Path = "Role/Prefabs/Pet_Dog", Scal = 1 },
            ["Pet_22"] = new ModelClass { Path = "Role/Prefabs/Pet_Cat", Scal = 1 },
            ["Pet_23"] = new ModelClass { Path = "Role/Prefabs/Pet_Mouse", Scal = 1 },
            ["Pet_24"] = new ModelClass { Path = "Role/Prefabs/Pet_Llama", Scal = 1 },
            ["Pet_25"] = new ModelClass { Path = "Role/Prefabs/Pet_Rabbit", Scal = 1 },
            ["Pet_26"] = new ModelClass { Path = "Role/Prefabs/Pet_Bear", Scal = 1 },
            ["Pet_27"] = new ModelClass { Path = "Role/Prefabs/Pet_Deer", Scal = 1 },
            ["Pet_28"] = new ModelClass { Path = "Role/Prefabs/Pet_Fox", Scal = 1 },
            ["Pet_29"] = new ModelClass { Path = "Role/Prefabs/Pet_Kangaroo", Scal = 1 },
            ["Pet_30"] = new ModelClass { Path = "Role/Prefabs/Pet_Pig", Scal = 1 },
            ["Pet_31"] = new ModelClass { Path = "Role/Prefabs/Pet_Turtle", Scal = 1 },
            ["Pet_51"] = new ModelClass { Path = "Role/Prefabs/Pet_Tiger", Scal = 1 },
            ["Pet_52"] = new ModelClass { Path = "Role/Prefabs/Pet_Long", Scal = 1 },
            ["Pet_53"] = new ModelClass { Path = "Role/Prefabs/Pet_QiLin", Scal = 1 },

        };

        private static string BaseResousePath = "Assets/HotUpdateResources/Main/";
        public static string GetMonsterModelByID(string key)
        {
            return BaseResousePath + MonsterDic[key].Path + ".prefab";
        }

        public static float GetMonsterScalModelByID(string key)
        {
            return MonsterDic[key].Scal;
        }

        private static Dictionary<string, string> EffDic = new Dictionary<string, string>
        {
            ["Attack_BJ"] = "Other/Effects/UI_eff/Attack_BJ",
            ["Attack_Miss"] = "Other/Effects/UI_eff/Attack_Miss",

            ["Eff_Ball2"] = "Other/Effects/UI_eff/eff__ui_zh_qiu_01",
            ["Eff_Ball3"] = "Other/Effects/UI_eff/eff__ui_zh_qiu_02",
            ["Eff_Ball4"] = "Other/Effects/UI_eff/eff__ui_zh_qiu_03",
            ["Eff_Ball6"] = "Other/Effects/UI_eff/eff__ui_zh_qiu_04",
            ["Eff_Ball7"] = "Other/Effects/UI_eff/eff__ui_zh_qiu_05",
            ["Eff_Ball8"] = "Other/Effects/UI_eff/eff__ui_zh_qiu_06",

            ["Boss_Romate_1_Bullet"] = "Role/MonsterPrefabs/Fight_Boss_Romate_1_Bullet",
            ["Boss_Romate_2_Bullet"] = "Role/MonsterPrefabs/Fight_Boss_Romate_2_Bullet",

            //人物技能
            ["eff_ZSAttack"] = "Other/Effects/spine_tx/Prefabs/Zs_Attack",
            ["eff_YHZ"] = "Other/Effects/spine_tx/Prefabs/Zs_YHZ",
            ["eff_SFZ"] = "Other/Effects/spine_tx/Prefabs/Zs_SFZ",
            ["eff_PJZ"] = "Other/Effects/spine_tx/Prefabs/Zs_PJZ",

            ["eff_FSAttack"] = "Other/Effects/spine_tx/Prefabs/FS_Attack",
            ["eff_JGMC"] = "Other/Effects/spine_tx/Prefabs/FS_JGMC",
            ["eff_MFFD"] = "Other/Effects/spine_tx/Prefabs/FS_MFFD",
            ["eff_MFFD_End"] = "Other/Effects/spine_tx/Prefabs/FS_MFFD_End",
            ["eff_XBS"] = "Other/Effects/spine_tx/Prefabs/FS_XBS",
            ["eff_XBS_Star"] = "Other/Effects/spine_tx/Prefabs/FS_XBS_Star",

            ["eff_JSAttack"] = "Other/Effects/spine_tx/Prefabs/JS_Attack",
            ["eff_LHFZ"] = "Other/Effects/spine_tx/Prefabs/JS_LHFZ",
            ["eff_LHFZ_End"] = "Other/Effects/spine_tx/Prefabs/JS_LHFZ_End",
            ["eff_LHFZ_Star"] = "Other/Effects/spine_tx/Prefabs/JS_LHFZ_Star",
            ["eff_SSZJS"] = "Other/Effects/spine_tx/Prefabs/JS_SSZJS",
            ["eff_YMS"] = "Other/Effects/spine_tx/Prefabs/JS_YMS",
            ["eff_ZDB"] = "Other/Effects/spine_tx/Prefabs/JS_ZDB",

            ["eff_JDQL"] = "Other/Effects/spine_tx/Prefabs/FS_JDQL",
            ["eff_SiXueZiZou"] = "Other/Effects/spine_tx/Prefabs/FS_SiXueZiZou",
            ["eff_MoFaJieJie"] = "Other/Effects/spine_tx/Prefabs/FS_MoFaJieJie",
            ["eff_MoLiQuZu"] = "Other/Effects/spine_tx/Prefabs/FS_MoLiQuZu", //魔力驱逐

            ["eff_DSS_Star"] = "Other/Effects/spine_tx/Prefabs/JS_DSS_Star",//毒砂术
            ["eff_DGS_Star"] = "Other/Effects/spine_tx/Prefabs/JS_DGS_Star",
            ["eff_DSS_Buffer"] = "Other/Effects/spine_tx/Prefabs/JS_DSS_Buffer",//毒砂术Buffer

            ["eff_QunTiZiYu"] = "Other/Effects/spine_tx/Prefabs/JS_QunTiZiYu",//群体治愈术
            ["eff_QunTiZiYu_Buffer"] = "Other/Effects/spine_tx/Prefabs/JS_QunTiZiYu_Buffer",//群体治愈术 Buffer

            ["eff_QunTiJingHua"] = "Other/Effects/spine_tx/Prefabs/JS_QunTiJingHua",//群体净化术
            ["eff_QunTiJingHua_Buffer"] = "Other/Effects/spine_tx/Prefabs/JS_QunTiJingHua_Buffer",//群体净化术 Buffer

            ["eff_SenLiJingTong"] = "Other/Effects/spine_tx/Prefabs/JS_SenLiJingTong",//神力精通
            ["eff_MoLiJingTong"] = "Other/Effects/spine_tx/Prefabs/JS_MoLiJingTong",//魔力精通
            ["eff_YongQiZuFu"] = "Other/Effects/spine_tx/Prefabs/JS_YongQiZuFu",//勇气祝福、

            ["eff_HuiGuangFanZao"] = "Other/Effects/spine_tx/Prefabs/JS_HuiGuangFanZao",//回光反照

            ["eff_MoFaZengFu"] = "Other/Effects/spine_tx/Prefabs/FS_MoFaZengFu", //魔法增幅
            ["eff_WuXieKeJi"] = "Other/Effects/spine_tx/Prefabs/FS_WXKJ", //无懈可击
            //["eff_XYZJ"] = "Other/Effects/prefabs/01Warrior_male/eff_Warrior_male_xunyingzhuiji",
            ["eff_XYZJ"] = "Other/Effects/spine_tx/Prefabs/Zs_XunYingZuiJi",//迅影追击

            //["eff_YJGX"] = "Other/Effects/prefabs/01Warrior_male/eff_Warrior_male_yijuegaoxia",
            ["eff_YJGX"] = "Other/Effects/spine_tx/Prefabs/Zs_YiJueGaoXia",//一决高下

            ["eff_XuanYunMengJi"] = "Other/Effects/spine_tx/Prefabs/Zs_XuanYunMengJi_Star",

            ["eff_SJCK"] = "Other/Effects/spine_tx/Prefabs/FS_SJCK",//

            ["eff_ChuanSong"] = "Other/Effects/spine_tx/Prefabs/eff_ChuanSong",//传送

            ["eff__WenYiZiRen"] = "Other/Effects/spine_tx/Prefabs/Zs_WenYiZiRen_Star",//瘟疫之刃

            ["Skill_ChaoFeng"] = "Other/Effects/spine_tx/Prefabs/Skill_ChaoFeng",
            ["Skill_KongJu"] = "Other/Effects/spine_tx/Prefabs/Skill_KongJu",
            ["Skill_ZhiMang"] = "Other/Effects/spine_tx/Prefabs/Skill_ZhiMang",
            ["Skill_JingHua"] = "Other/Effects/spine_tx/Prefabs/Skill_JingHua",

            ["Buffer_32"] = "Other/Effects/spine_tx/Prefabs/FS_TSBB",

            ["Buffer_9"] = "Other/Effects/spine_tx/Prefabs/Buffer_9",//
            ["Buffer_10"] = "Other/Effects/spine_tx/Prefabs/Buffer_10",
            ["Buffer_11"] = "Other/Effects/spine_tx/Prefabs/Buffer_11",
            ["Buffer_12"] = "Other/Effects/spine_tx/Prefabs/Buffer_12",
            ["Buffer_13"] = "Other/Effects/spine_tx/Prefabs/Buffer_13",
            ["Buffer_22"] = "Other/Effects/spine_tx/Prefabs/Buffer_22",
            ["Buffer_24"] = "Other/Effects/spine_tx/Prefabs/Buffer_24",
            ["Buffer_25"] = "Other/Effects/spine_tx/Prefabs/Buffer_25",
            ["Buffer_26"] = "Other/Effects/spine_tx/Prefabs/Buffer_26",
            ["Buffer_31"] = "Other/Effects/spine_tx/Prefabs/Buffer_31",
            ["Buffer_33"] = "Other/Effects/spine_tx/Prefabs/Buffer_33",
            ["Buffer_39"] = "Other/Effects/spine_tx/Prefabs/Buffer_39",
            ["Buffer_42"] = "Other/Effects/spine_tx/Prefabs/Buffer_42",
            ["Buffer_43"] = "Other/Effects/spine_tx/Prefabs/Buffer_43",
            ["Buffer_44"] = "Other/Effects/spine_tx/Prefabs/Buffer_44",
            ["Buffer_54"] = "Other/Effects/spine_tx/Prefabs/Buffer_54",
            ["Buffer_Feng"] = "Role/MonsterPrefabs/Fight_Boss_8_Feng",
            ["Buffer_yemanren"] = "Other/Effects/spine_tx/Prefabs/Buffer_YeManRen",
            ["Buffer_ZhangYu"] = "Other/Effects/spine_tx/Prefabs/Buffer_ZhangYu",
            ["Buffer_Boss_Du"] = "Other/Effects/spine_tx/Prefabs/Buffer_Boss_Du",
            ["Buffer_PingHeng"] = "Other/Effects/spine_tx/Prefabs/Buffer_PingHeng",

            ["Debuff1"] = "Other/Effects/spine_tx/Prefabs/DeBuffer_1",
            ["Debuff2"] = "Other/Effects/spine_tx/Prefabs/DeBuffer_2",
            ["Debuff3"] = "Other/Effects/spine_tx/Prefabs/DeBuffer_3",
            ["Debuff4"] = "Other/Effects/spine_tx/Prefabs/DeBuffer_4",

            //怪物
            //["bullet_orange"] = "Other/Effects/prefabs/03monster/eff_monster_attack_03",
            //["bullet_pink"] = "Other/Effects/prefabs/03monster/eff_monster_attack_02",
            //["bullet_blue"] = "Other/Effects/prefabs/03monster/eff_monster_attack_01",
            //["eff_MonsterDeath"] = "Other/Effects/prefabs/03monster/eff_monster_die",
            //["eff_MonsterBorn"] = "Other/Effects/prefabs/03monster/eff_monster_appear",

            //["Hurt_1"] = "Other/Effects/prefabs/03monster/eff_monster_hit_male",
            //["Hurt_2"] = "Other/Effects/prefabs/03monster/eff_monster_hit_female",
            ["Hurt_1"] = "Other/Effects/spine_tx/Prefabs/Zs_Hurt",
            ["Hurt_2"] = "Other/Effects/spine_tx/Prefabs/FS_Hurt",
            ["Hurt_3"] = "Other/Effects/spine_tx/Prefabs/JS_Hurt",
            ["Monster_zd"] = "Other/Effects/UI_eff/monster_zd",

            ["Romate_7_Bullet"] = "Role/MonsterPrefabs/Fight_Boss_Romate_7_Bullet",
            ["Romate_6_Bullet"] = "Role/MonsterPrefabs/Fight_Boss_Romate_6_Bullet",
            ["Romate_6_Hit"] = "Role/MonsterPrefabs/Fight_Boss_Romate_6_Hit",
            ["Romate_6_Skill"] = "Role/MonsterPrefabs/Fight_Boss_Romate_6_Skill",

            ["BattleMap_1"] = "Terrain/Battle_1",
            ["BattleMap_2"] = "Terrain/Battle_2",

            ["LvUp"] = "Other/Effects/spine_tx/Prefabs/Fight_LvLp",
            ["UI_LvUP"] = "Other/Effects/UI_eff/UI_LvUP",

            //掉落特效
            //["drop_3"] = "Other/Effects/prefabs/04other/eff_other_common_diaoluo_02",
            //["drop_4"] = "Other/Effects/prefabs/04other/eff_other_common_diaoluo_03",
            //["drop_5"] = "Other/Effects/prefabs/04other/eff_other_common_diaoluo_04",
            //["drop_6"] = "Other/Effects/prefabs/04other/eff_other_common_diaoluo_05",
            //["drop_7"] = "Other/Effects/prefabs/04other/eff_other_common_diaoluo_06",

            ["UIFly"] = "Other/Effects/UI_eff/eff_other_ui_fly",
            ["UIBoom"] = "Other/Effects/UI_eff/eff_other_ui_boom",

            ["UIPartnerCall"] = "Other/Effects/spine_tx/Ui_huoban/spine",
            ["UIKuang"] = "Other/Effects/spine_tx/Ui_kuang/spine",

            ["UIAttack_1"] = "Other/Effects/UI_eff/UI_ZS_Attack",
            ["UIAttack_2"] = "Other/Effects/UI_eff/UI_FS_Attack",
            ["UIAttack_3"] = "Other/Effects/UI_eff/UI_JS_Attack",

            ["UI_Hurt"] = "Other/Effects/UI_eff/UI_Hurt",
            ["UI_ExpUp"] = "Other/Effects/UI_eff/eff__ui_up",

            // 主战场特效
            // 战士
            ["UI_Hall_ZS_Skill1_BanYueZhan"] = "Other/Effects/GameSpine_A/UI_Hall_ZS_Skill1_BanYueZhan",
            ["UI_Hall_ZS_Skill1_Hit"] = "Other/Effects/GameSpine_A/UI_Hall_ZS_Skill1_Hit",
            ["UI_Hall_ZS_Skill3_LieYanZhan"] = "Other/Effects/GameSpine_A/UI_Hall_ZS_Skill3_LieYanZhan",
            ["UI_Hall_ZS_Skill3_Hit"] = "Other/Effects/GameSpine_A/UI_Hall_ZS_Skill3_Hit",
            ["UI_Hall_ZS_Skill2_ShiZiZhan"] = "Other/Effects/GameSpine_A/UI_Hall_ZS_Skill2_ShiZiZhan",

            // 法师
            ["UI_Hall_FS_Skill1_GuangQiu"] = "Other/Effects/GameSpine_A/UI_Hall_FS_Skill1_GuangQiu",
            ["UI_Hall_FS_Skill1_Hit"] = "Other/Effects/GameSpine_A/UI_Hall_FS_Skill1_Hit",
            ["UI_Hall_FS_Skill2_Hit"] = "Other/Effects/GameSpine_A/UI_Hall_FS_Skill2_Hit",
            ["UI_Hall_FS_Skill3_MaiChongBo"] = "Other/Effects/GameSpine_A/UI_Hall_FS_Skill3_MaiChongBo",
            ["UI_Hall_FS_Skill3_Hit"] = "Other/Effects/GameSpine_A/UI_Hall_FS_Skill3_Hit",

            // 祭祀
            ["UI_Hall_JS_Skill1_Cast"] = "Other/Effects/GameSpine_A/UI_Hall_JS_Skill1_Cast",
            ["UI_Hall_JS_Skill1_Hit"] = "Other/Effects/GameSpine_A/UI_Hall_JS_Skill1_Hit",
            ["UI_Hall_JS_Skill2_Di"] = "Other/Effects/GameSpine_A/UI_Hall_JS_Skill2_Di",
            ["UI_Hall_JS_Skill2_Hit"] = "Other/Effects/GameSpine_A/UI_Hall_JS_Skill2_Hit",
            ["UI_Hall_JS_Skill3_Hit"] = "Other/Effects/GameSpine_A/UI_Hall_JS_Skill3_Hit",

            // 大招
            // 荆棘牢笼
            ["UI_Hall_Common_Skill_JingJiLaoLong"] = "Other/Effects/GameSpine_A/UI_Hall_Common_Skill_JingJiLaoLong",
            ["UI_Hall_Common_Skill_ShiFa"] = "Other/Effects/GameSpine_A/UI_Hall_Common_Skill_ShiFa",

            // 女神降临
            ["UI_Hall_Common_Skill_NvShenJiangLin"] = "Other/Effects/GameSpine_A/UI_Hall_Common_Skill_NvShenJiangLin",
            ["UI_Hall_Common_Skill_HuiFu"] = "Other/Effects/GameSpine_A/UI_Hall_Common_Skill_HuiFu",

            // 时空之轮
            ["UI_Hall_Common_Skill_ShiKongZhiLun"] = "Other/Effects/GameSpine_A/UI_Hall_Common_Skill_ShiKongZhiLun",
            ["UI_Hall_Common_Skill_ShiKongZhiLun1"] = "Other/Effects/GameSpine_A/UI_Hall_Common_Skill_ShiKongZhiLun1",
            ["UI_Hall_Common_Skill_ShiKongZhiLunHit"] = "Other/Effects/GameSpine_A/UI_Hall_Common_Skill_ShiKongZhiLunHit",

            ////// 
            ["UI_Hall_Common1_Skill_ShiFa_Zong"] = "Other/Effects/GameSpine_A/UI_Hall_Common1_Skill_ShiFa_Zong",
            // 虚空坍缩
            ["UI_Hall_Common1_Skill1_Down"] = "Other/Effects/GameSpine_A/UI_Hall_Common1_Skill1_Down",
            ["UI_Hall_Common1_Skill1_Hit"] = "Other/Effects/GameSpine_A/UI_Hall_Common1_Skill1_Hit",
            ["UI_Hall_Common1_Skill1_Up"] = "Other/Effects/GameSpine_A/UI_Hall_Common1_Skill1_Up",

            // 虚空抓
            ["UI_Hall_Common1_Skill2_Down"] = "Other/Effects/GameSpine_A/UI_Hall_Common1_Skill2_Down",
            ["UI_Hall_Common1_Skill2_Hit"] = "Other/Effects/GameSpine_A/UI_Hall_Common1_Skill2_Hit",
            ["UI_Hall_Common1_Skill2_Up"] = "Other/Effects/GameSpine_A/UI_Hall_Common1_Skill2_Up",

            // 水龙斩
            ["UI_Hall_Common1_Skill3_Down"] = "Other/Effects/GameSpine_A/UI_Hall_Common1_Skill3_Down",
            ["UI_Hall_Common1_Skill3_Hit"] = "Other/Effects/GameSpine_A/UI_Hall_Common1_Skill3_Hit",
            ["UI_Hall_Common1_Skill3_Up"] = "Other/Effects/GameSpine_A/UI_Hall_Common1_Skill3_Up",

            // 漫天箭雨
            ["UI_Hall_Common1_Skill4_LuoJian"] = "Other/Effects/GameSpine_A/UI_Hall_Common1_Skill4_LuoJian",
            ["UI_Hall_Common1_Skill4_DiLie"] = "Other/Effects/GameSpine_A/UI_Hall_Common1_Skill4_DiLie",
            ["UI_Hall_Common1_Skill4_BaoZha"] = "Other/Effects/GameSpine_A/UI_Hall_Common1_Skill4_BaoZha",

            // 忍者碎星击
            ["UI_Hall_Common1_Skill5_ShiFa"] = "Other/Effects/GameSpine_A/UI_Hall_Common1_Skill5_ShiFa",
            ["UI_Hall_Common1_Skill5_HuanYing"] = "Other/Effects/GameSpine_A/UI_Hall_Common1_Skill5_HuanYing",
            ["UI_Hall_Common1_Skill5_DiLie"] = "Other/Effects/GameSpine_A/UI_Hall_Common1_Skill5_DiLie",
            ["UI_Hall_Common1_Skill5_BaoZha"] = "Other/Effects/GameSpine_A/UI_Hall_Common1_Skill5_BaoZha",

            // 碎星地裂斩
            ["UI_Hall_Common2_Skill1_ShiFa"] = "Other/Effects/GameSpine_A/UI_Hall_Common2_Skill1_ShiFa",
            ["UI_Hall_Common2_Skill1_Up"] = "Other/Effects/GameSpine_A/UI_Hall_Common2_Skill1_Up",
            ["UI_Hall_Common2_Skill1_ShouJi"] = "Other/Effects/GameSpine_A/UI_Hall_Common2_Skill1_ShouJi",
            ["UI_Hall_Common2_Skill1_Down"] = "Other/Effects/GameSpine_A/UI_Hall_Common2_Skill1_Down",
        };

        public static string GetEffPath(string key)
        {
            return BaseResousePath + EffDic[key] + ".prefab";
        }

        public static bool CheakEff(string key)
        {
            return EffDic.ContainsKey(key);
        }

        public const string GameSettingKay = "PlayerPrefs_GameSettingKay";

        public enum PlayerPrefsKey
        {
            BagRecycle,
            BoxRecycle,
            PetBuy,
            HangUpCDTime,
        }

        //长按按钮响应时间
        public static readonly float TouchTime = 0.5f;

        //角色头像图集
        public const string RoleHeadAtlas = "Common/Atlas/RoleImg";
        public const string CommonAtlas = "Common/Atlas/Common";
        public const string CommonBGAtlas = "Common/Atlas/CommonBG";
        public const string FirstPayAtlas = "Common/Atlas/FirstPay";
        public const string GiftBagBanner = "Common/Atlas/GiftBagBanner";
        public const string FightAtlas = "Common/Atlas/Fight";
        public const string CreatRoleAtlas = "Common/Atlas/CreatRole";
        public const string WarFightAtlas = "Common/Atlas/WarFight";
        public const string ChatAtlas = "Common/Atlas/Chat";
        public const string PartnerCall = "Common/Atlas/PartnerCall";
        public const string CompetitionEntranceAtlas = "Common/Atlas/CompetitionEntrance";
        public const string ActivityEntranceAtlas = "Common/Atlas/ActivityEntrance";
        public const string CardAtlas = "Common/Atlas/Card";
        public const string HallAtlas = "Common/Atlas/Hall";
        public const string HallBattleAtlas = "Common/Atlas/HallBattle";

        public const string RoleIcon = "Common/Atlas/Head";

        public static string[] RelicItemStrs = { "限定时间内击杀怪物", "通关时总血量高于50%", "通关时所有角色均存活" };

        public static string LoginPath = "Main/Scene/Login";
        public static string HallPath = "Main/Scene/Hall";
        public static string NorFightPath = "Main/Scene/pve_scene";
        public static string TurnBasedFightPath = "Main/Scene/turn_based_scene";

        // 文本置灰颜色
        public static Color GrayTextColor => new Color(0.263f, 0.263f, 0.263f);
        // 是否首次打开小游戏？
        public const string IsFirstOpenMiniGame = "IsFirstOpenMiniGame";
        // 头像品质资源前缀
        public const string HeadIconQualityPrefix = "bg_avatar_icon_";
    }
}
