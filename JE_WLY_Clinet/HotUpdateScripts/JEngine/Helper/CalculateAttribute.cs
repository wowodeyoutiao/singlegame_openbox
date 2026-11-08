using BestHTTP.Examples;
using BestHTTP.SecureProtocol.Org.BouncyCastle.Asn1.Ocsp;
using HallBattle;
using HotUpdateScripts;
using HotUpdateScripts.Proto;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using JEngine.Singleton;
using pb;
using Spine;
using System;
using System.Collections;
using System.Collections.Generic;
using System.Configuration;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Threading;
using Unity.Jobs;
using UnityEngine;
using static CFGItem_Fate;
using static UIFightFormation;
using static UnityEngine.Networking.UnityWebRequest;

namespace JEngine.Helper
{
    /// <summary>
    /// 本类用于计算角色基础属性
    /// </summary>
    public class CalculateAttribute : SingletonTemplate<CalculateAttribute>
    {
        //存储需要乘以系数的属性及系数信息
        public readonly Dictionary<pb.Attribute, pb.Attribute> AttrRelationShip = new() {
            {pb.Attribute.PHYSICAL_ATTACK, pb.Attribute.PHYSICAL_ATTACK_PERCENTAGE },
            {pb.Attribute.MAGIC_ATTACK, pb.Attribute.MAGIC_ATTACK },
            {pb.Attribute.PHYSICAL_DEFENSE, pb.Attribute.PHYSICAL_DEFENSE_PERCENTAGE },
            {pb.Attribute.MAGIC_DEFENSE, pb.Attribute.MAGIC_DEFENSE_PERCENTAGE },
            {pb.Attribute.DODGE, pb.Attribute.DODGE_PERCENTAGE },
            {pb.Attribute.HIT, pb.Attribute.HIT_PERCENTAGE },
            {pb.Attribute.CRITICAL_STRIKE, pb.Attribute.CRITICAL_STRIKE_PERCENTAGE },
            {pb.Attribute.TOUGHNESS, pb.Attribute.TOUGHNESS_PERCENTAGE },
            {pb.Attribute.MAX_HEALTH, pb.Attribute.HEALTH_PERCENTAGE },
            {pb.Attribute.ATTACK_SPEED, pb.Attribute.ATTACK_SPEED_PERCENTAGE },
            {pb.Attribute.MOVE_SPEED, pb.Attribute.MOVE_SPEED_PERCENTAGE },
        };

        //存储需要乘以系数的属性
        public readonly HashSet<pb.Attribute> AdjustAttr = new() {
            pb.Attribute.PHYSICAL_ATTACK,
            pb.Attribute.MAGIC_ATTACK,
            pb.Attribute.PHYSICAL_DEFENSE,
            pb.Attribute.MAGIC_DEFENSE,
            pb.Attribute.DODGE,
            pb.Attribute.HIT,
            pb.Attribute.CRITICAL_STRIKE,
            pb.Attribute.TOUGHNESS,
            pb.Attribute.MAX_HEALTH,
            pb.Attribute.ATTACK_SPEED,
            pb.Attribute.MOVE_SPEED,
        };

        public readonly double BaseRate = 10000f; //万分比分母

        #region 获取穿上装备的战力
        public int GetPowerWhenPutOnEquipment(int job, int teamIndex, pb.ItemInfo tryPutOnItemInfo)
        {
            (bool isDoNextStep, pb.ItemInfo putOnItemInfo, int power) = checkIsPuton(tryPutOnItemInfo, teamIndex);

            //如果该装备已经穿上就不需要进行下一步计算，直接返回战力信息
            if (!isDoNextStep) return power;

            List<Stats> statsPanelList = teamIndex == 0 ? RoleM.Instance.GetRoleAllStats() : PartnerM.Instance.GetOnePartnerStats(PartnerM.Instance.GetOnePartnerSlotInfo(teamIndex).PartnerUid).Attributes;

            List<Stats> newStatsList = TryPutOnAttributeData(statsPanelList, putOnItemInfo, tryPutOnItemInfo,job,teamIndex);

            var skillList = teamIndex == 0 ? SkillM.Instance.GetSkillInfoList() : PartnerM.Instance.GetSkillInfoByPartnerId(PartnerM.Instance.GetOnePartnerSlotInfo(teamIndex).PartnerUid);
            return CalculatePower(newStatsList, skillList, job);
        }

        public void Test()
        {
            List<Stats> statsList = RoleM.Instance.GetRoleAllStats();
            var skillList = SkillM.Instance.GetSkillInfoList();

            //用于返回还原的结果
            List<Stats> result = new List<Stats>();
            //记录字典信息
            Dictionary<pb.Attribute, Stats> dicAttr = new Dictionary<pb.Attribute, Stats>();

            convertBaseStats(statsList, result, dicAttr);

            foreach (var item in statsList)
            {
                var stats = new Stats() { attr = item.attr, Num = item.Num };

                if (!dicAttr.ContainsKey(item.attr))
                    dicAttr.Add(item.attr, stats);

                result.Add(stats);
            }
            
            debugInfo("原始数据  =====", dicAttr);

            var Power = CalculatePower(statsList, skillList, 1);
            string strInfo = "面板属性。。。。战力："+ Power + "\n";
            statsList.Sort((a, b) => (int)a.attr - (int)b.attr);
            foreach (var item in statsList)
            {
                strInfo += $"{item.attr.ToString()}_{item.Num}\n";
            }

            UnityEngine.Debug.LogWarning(strInfo);


        }

        //转换成面板属性
        private void convertPanelStats(List<Stats> statsList)
        {
            
        }

        //检查该装备是否穿上
        private (bool, pb.ItemInfo, int) checkIsPuton(pb.ItemInfo itemInfo, int teamIndex)
        {
            var equipCFG = CFG_Item.GetById(itemInfo.CfgId);
            int power = 0;
            pb.ItemInfo putOnEquipment = null;

            teamIndex = teamIndex == -1 ? 0 : teamIndex;
            if (teamIndex == 0)
            {
                putOnEquipment = RoleM.Instance.GetSlotInfoBySolt((pb.EquipSlot)equipCFG.Equiptype).Item;
                power = RoleM.Instance.Combat;
            }
            else
            {
                var solt = PartnerM.Instance.GetOnePartnerSlotInfo(teamIndex);
                if (solt.PartnerUid != 0)
                {
                    putOnEquipment = PartnerM.Instance.GetOnePartnerOneEquip(teamIndex, (int)equipCFG.Equiptype);

                    power = PartnerM.Instance.GetOnePartnerStats(solt.PartnerUid).Combat;
                }
            }

            //是否进行下一步计算
            bool isDoNextStep = putOnEquipment == null ? false : itemInfo.Uid != putOnEquipment.Uid;

            UnityEngine.Debug.LogWarning($"尝试装备ID{itemInfo.Uid}_已装备ID{(putOnEquipment== null?-1:putOnEquipment.Uid)}");
            return (isDoNextStep, putOnEquipment, power);
        }
        #endregion

        #region 战力计算
        /// <summary>
        /// 战力计算
        /// </summary>
        /// <param name="StatsList"></param>
        /// <param name="skillList"></param>
        /// <param name="job"></param>
        /// <returns></returns>
        public int CalculatePower(List<Stats> StatsList, List<SkillInfo> skillList, int job)
        {
            Dictionary<pb.Attribute, int> dicStats = getStatsDic(StatsList);

            float attackValue = getAttackValue(dicStats, job); //职业攻击加成

            float skillBonus = getSkillCombat(skillList); //技能职业加成

            int maxHealth = tryGetDicVal(dicStats, pb.Attribute.MAX_HEALTH);     // 面板生命值 12

            int physicalDefense = tryGetDicVal(dicStats, pb.Attribute.PHYSICAL_DEFENSE); // 面板物防值	4

            int magicDefense = tryGetDicVal(dicStats, pb.Attribute.MAGIC_DEFENSE);    // 面板魔防值	5

            int hitValue = tryGetDicVal(dicStats, pb.Attribute.HIT);  // 命中值	7

            int dodgeValue = tryGetDicVal(dicStats, pb.Attribute.DODGE); // 闪避值	6

            int criticalValue = tryGetDicVal(dicStats, pb.Attribute.CRITICAL_STRIKE);  // 暴击值	8

            int toughnessValue = tryGetDicVal(dicStats, pb.Attribute.TOUGHNESS);   // 韧性值	9

            // 二级属性
            int lifeStealTriggerPct = 0; //吸血触发百分比   =============//未投放  

            int lifeStealPct = tryGetDicVal(dicStats, pb.Attribute.SUCKING_BLOOD_PERCENTAGE);  // 吸血比例	40

            int finalDamageReductionPct = tryGetDicVal(dicStats, pb.Attribute.DMT_FINAL_REDUCE_PERCENTAGE); // 最终减伤百分比	35

            int PvEDamageReductionPct = tryGetDicVal(dicStats, pb.Attribute.DMT_REDUCE_PERCENTAGE_MONSTER);// 对怪减伤百分比	33

            int PvPDamageReductionPct = tryGetDicVal(dicStats, pb.Attribute.DMT_REDUCE_PERCENTAGE_PLAYER); // 对人减伤百分比	32

            int criticalDamagePct = tryGetDicVal(dicStats, pb.Attribute.CRITICAL_DAMAGE_PERCENTAGE);// 暴击伤害增加百分比	24

            int reflectDamageTriggerPct = 0;  // 反伤触发百分比  =============//未投放  

            int reflectDamagePct = tryGetDicVal(dicStats, pb.Attribute.REFLECT_DAMAGE_PERCENTAGE);// 反伤比例	42

            int PvEDamageIncreasePct = tryGetDicVal(dicStats, pb.Attribute.PVE_DAMAGE_PERCENTAGE); // 对怪增伤百分比	26

            int PvPDamageIncreasePct = tryGetDicVal(dicStats, pb.Attribute.PVP_DAMAGE_PERCENTAGE);// 对人增伤百分比	25

            int finalDamageIncreasePct = tryGetDicVal(dicStats, pb.Attribute.DMT_FINAL_INCREASE_PERCENTAGE);// 最终增伤百分比	39

            // 三级属性
            int PvEDamageIncrease = tryGetDicVal(dicStats, pb.Attribute.DMT_INCREASE_FIXED_MONSTER);// 对怪增伤固定值	37

            int PvEDamageReduction = tryGetDicVal(dicStats, pb.Attribute.DMT_REDUCE_FIXED_MONSTER);// 对怪减伤固定值	31

            int PvPDamageReduction = tryGetDicVal(dicStats, pb.Attribute.DMT_REDUCE_FIXED_PLAYER);// 对人减伤固定值	30

            int PvPDamageIncrease = tryGetDicVal(dicStats, pb.Attribute.DMT_INCREASE_FIXED_PLAYER);// 对人增伤固定值	36

            int finalDamageReduction = tryGetDicVal(dicStats, pb.Attribute.DMT_FINAL_REDUCE_FIXED);// 最终减伤固定值	34

            int finalDamageIncrease = tryGetDicVal(dicStats, pb.Attribute.DMT_FINAL_INCREASE_FIXED);// 最终增伤固定值	38


            int deffenseCoefficient = CFG_XMLConfig.xmlConfig.BaseValue.Power.defvalue;// 防御系数

            int effectiveHPCoefficient = CFG_XMLConfig.xmlConfig.BaseValue.Power.lifevalue;// 有效生命系数

            int effectivePowerCoefficient = CFG_XMLConfig.xmlConfig.BaseValue.Power.attackvalue;// 有效输出系数

            int hitCoefficient = CFG_XMLConfig.xmlConfig.BaseValue.Power.hitvalue;// 命中系数

            int dodgeCoefficient = CFG_XMLConfig.xmlConfig.BaseValue.Power.dodgevalue;// 闪避系数

            int criticalCoefficient = CFG_XMLConfig.xmlConfig.BaseValue.Power.critvalue;// 暴击系数

            int toughnessCoefficient = CFG_XMLConfig.xmlConfig.BaseValue.Power.toughness;// 韧性系数


            // 有效生命值计算：(面板生命值×【1+吸血触发百分比×（1+吸血比例）】+防御系数×(（面板物防值+面板魔防值）)⁄2)/(（1-最终减伤百分比-对怪减伤百分比/2-对人减伤百分比/2）)+(（对怪减伤固定值+对人减伤固定值）)/2+最终减伤固定值
            double lifeWithBlood = maxHealth * (1 + (lifeStealTriggerPct / BaseRate) * (1 + lifeStealPct / BaseRate));

            double defenseValue = deffenseCoefficient * ((physicalDefense + magicDefense) / 2.0d);
            // 计算总减伤百分比，并确保不超过阈值
            double totalReductionPct = finalDamageReductionPct + PvEDamageReductionPct / 2.0d + PvPDamageReductionPct / 2.0f;
            // 减伤系数，确保不为0或负数
            double reductionCoeff = 1 - totalReductionPct / BaseRate;

            if (reductionCoeff <= 0)
                reductionCoeff = 0.01d; // 最小减伤系数，避免除零或负数

            //有效生命
            double effectiveHP = (lifeWithBlood + defenseValue) / reductionCoeff + (PvEDamageReduction + PvPDamageReduction) / 2.0 + finalDamageReduction;



            // 有效输出计算：面板主攻击值×【1+增伤百分比+对人增伤百分比/2+对怪增伤百分比/2+暴击伤害增加百分比+反伤触发百分比×（1+反伤比例）】×(1+技能加成效果/100)+(（对怪增伤固定值+对人增伤固定值）)/2+最终增伤固定值
            double damagePercentage = (finalDamageIncreasePct + PvEDamageIncreasePct / 2.0d + PvPDamageIncreasePct / 2.0d + criticalDamagePct) / BaseRate + (reflectDamageTriggerPct / BaseRate) * (1 + reflectDamagePct / BaseRate);

            double skillMultiplier = (1 + skillBonus / 100.0d);

            double fixedDamageIncrease = (PvEDamageIncrease + PvPDamageIncrease) / 2.0d + finalDamageIncrease;

            //有效输出
            double effectivePower = attackValue * (1 + damagePercentage) * skillMultiplier + fixedDamageIncrease;

            // 战力计算：√(有效生命×有效生命系数)×√(有效输出×有效输出系数)+命中值×命中系数+闪避值×闪避系数+暴击值×暴击系数+韧性值×韧性系数+击晕值×击晕系数+抗击晕值×抗击晕系数+中毒值×中毒系数+毒抗值×毒抗系数
            double combat = Math.Sqrt(effectiveHP * effectiveHPCoefficient) *
                Math.Sqrt(effectivePower * effectivePowerCoefficient) +
                hitValue * hitCoefficient +
                dodgeValue * dodgeCoefficient +
                criticalValue * criticalCoefficient +
                toughnessValue * toughnessCoefficient;

            return (int)Math.Round(combat);
        }

        private Dictionary<pb.Attribute, int> getStatsDic(List<Stats> statsList)
        {
            Dictionary<pb.Attribute, int> result = new();
            foreach (var item in statsList)
            {
                if (result.ContainsKey(item.attr))
                    result[item.attr] += item.Num;
                else
                    result.Add(item.attr, item.Num);
            }
            return result;
        }

        private int tryGetDicVal(Dictionary<pb.Attribute, int> result, pb.Attribute attr)
        {
            if (result == null) return 0;

            if (!result.ContainsKey(attr)) return 0;

            return result[attr];
        }

        //获取技能加成
        private float getSkillCombat(List<SkillInfo> skillList)
        {
            if (skillList == null) return 0;
            float skillCombat = 0;

            for (int i = 0; i < skillList.Count; i++)
            {
                var cfgSkillInfo = CFG_Skill.GetById(skillList[i].CfgId);
                skillCombat += cfgSkillInfo.Powerfactor * (1 + skillList[i].Advanced);
            }

            return skillCombat;
        }

        //获取职业加成
        private float getAttackValue(Dictionary<pb.Attribute, int> dicStats, int job)
        {
            //job的枚举值  1 战士；2 法师；3 祭祀；11 伙伴战士；12伙伴法师；13 伙伴祭祀；21 宠物；
            if (job % 10 == 1)
                return tryGetDicVal(dicStats, pb.Attribute.PHYSICAL_ATTACK);

            return tryGetDicVal(dicStats, pb.Attribute.MAGIC_ATTACK);
        }

        #endregion

        #region 尝试穿装备
        /// <summary>
        /// 尝试穿装备
        /// </summary>
        /// <param name="panelData">服务端原始数据</param>
        /// /// <param name="curItemInfo">当前的装备数据</param>
        /// <param name="tryPutOnItemInfo">预装备数据</param>
        /// <returns></returns>
        public List<Stats> TryPutOnAttributeData(List<Stats> panelData, pb.ItemInfo curItemInfo, pb.ItemInfo tryPutOnItemInfo,int job = 0,int teamIndex = 0)
        {

            /***********************************已有结论，面板属性有百分比的相关属性除以百分比就是基础属性*************************************/



            //用于返回还原的结果
            List<Stats> result = new List<Stats>();
            //记录字典信息
            Dictionary<pb.Attribute, Stats> dicAttr = new Dictionary<pb.Attribute, Stats>();

            convertBaseStats(panelData, result, dicAttr);

            //foreach (var item in originalData)
            //{
            //    var stats = new Stats() { attr = item.attr, Num = item.Num };

            //    if (!dicAttr.ContainsKey(item.attr))
            //        dicAttr.Add(item.attr, stats);

            //    result.Add(stats);
            //}

            var skillList = teamIndex == 0 ? SkillM.Instance.GetSkillInfoList() : PartnerM.Instance.GetSkillInfoByPartnerId(PartnerM.Instance.GetOnePartnerSlotInfo(teamIndex).PartnerUid);

            
            #region  Test Code
            //debugPanelDataConvert2BaseData(panelData);
            //TestCode(skillList, job);
            #endregion


            var Power = CalculatePower(result, skillList, job);
            debugInfo("原始数据  战力：" + Power, dicAttr);

            //装备当前装备属性
            Dictionary<pb.Attribute, int> combinItemAttr = getItemAttributes(curItemInfo);

            //遍历装备属性并移除当前装备所加的属性
            removeEquipAttrFromStatsDic(combinItemAttr, dicAttr,curItemInfo);
            Power = CalculatePower(result, skillList, job);
            debugInfo("脱下当前装备后的属性 战力："+ Power, dicAttr);


            //获取新装备属性
            Dictionary<pb.Attribute, int> newCombinItemAttr = getItemAttributes(tryPutOnItemInfo);
            //遍历装备属性并增加当前装备所加的属性
            addAttrToStatsDic(newCombinItemAttr, dicAttr);
            Power = CalculatePower(result, skillList, job);
            debugInfo("穿上新装备后的属性 Power" + Power, dicAttr);

            return result;
        }

        private void TestCode(List<SkillInfo> skillList, int job = 0)
        {
            //var test = ProtoC.BaseAttributes;

            //List<Stats> result = new List<Stats>();

            //Dictionary<pb.Attribute, Stats> dicAttr = new Dictionary<pb.Attribute, Stats>();

            //Dictionary<pb.Attribute, int> temp = new();
            //foreach (var item in test)
            //{
            //    //if(!temp.ContainsKey(item.attr))
            //    try {
            //        temp.Add(item.attr, item.Num);
            //    }
            //    catch {
            //        UnityEngine.Debug.LogWarning("有相同项");
            //    }

                    
            //    UnityEngine.Debug.LogWarning(item.attr.ToString() + "___________" + item.Num);
            //}

            ////为了方便装备新增属性不在服务器下发的属性列表中，此处构造一个全部的属性列表
            //foreach (pb.Attribute attr in Enum.GetValues(typeof(pb.Attribute)))
            //{
            //    if (attr == pb.Attribute.INVALID_Attribute || attr == pb.Attribute.MAX_Attribute) continue;

            //    var stats = new Stats() { attr = attr, Num = 0 };
            //    result.Add(stats);

            //    //当前属性服务器已有数据，则直接赋值
            //    if (temp.ContainsKey(attr) )
            //    {
            //        if (AdjustAttr.Contains(attr))
            //        {
            //            var percent = AttrRelationShip[attr];
            //            if (temp.ContainsKey(percent))
            //            {
            //                stats.Num = (int)(temp[attr] * (1 + temp[percent] / BaseRate));
            //            }
            //            else
            //                stats.Num = temp[attr];
            //        }
            //        else
            //            stats.Num = temp[attr];
            //    }

            //    dicAttr.Add(attr, stats);
            //}

            //var Power = CalculatePower(result, skillList, job);
            //debugInfo("服务端原始战力转面板后计算的战力：" + Power, dicAttr);
        }

        private void debugPanelDataConvert2BaseData(List<Stats> panelData)
        {
            //List<Stats> result1 = new List<Stats>();
            //Dictionary<pb.Attribute, Stats> dicAttr1 = new Dictionary<pb.Attribute, Stats>();
            //convertToBaseStats(panelData, result1, dicAttr1);

            //debugInfo("面板数据转基础数据：", dicAttr1);



            //var test = ProtoC.BaseAttributes;

            //List<Stats> result = new List<Stats>();

            //Dictionary<pb.Attribute, Stats> dicAttr = new Dictionary<pb.Attribute, Stats>();

            //Dictionary<pb.Attribute, int> temp = new();
            //foreach (var item in test)
            //{
            //    //if(!temp.ContainsKey(item.attr))
            //    try
            //    {
            //        temp.Add(item.attr, item.Num);
            //    }
            //    catch
            //    {
            //        UnityEngine.Debug.LogWarning("有相同项");
            //    }
            //}

            ////为了方便装备新增属性不在服务器下发的属性列表中，此处构造一个全部的属性列表
            //foreach (pb.Attribute attr in Enum.GetValues(typeof(pb.Attribute)))
            //{
            //    if (attr == pb.Attribute.INVALID_Attribute || attr == pb.Attribute.MAX_Attribute) continue;

            //    var stats = new Stats() { attr = attr, Num = 0 };
            //    result.Add(stats);

            //    //当前属性服务器已有数据，则直接赋值
            //    if (temp.ContainsKey(attr))
            //    {
            //        stats.Num = temp[attr];
            //    }

            //    dicAttr.Add(attr, stats);
            //}

            //debugInfo("服务端原始基础信息：" , dicAttr);

        }


        private void debugInfo(string titleInfo,Dictionary<pb.Attribute, Stats> dicAttr)
        {
            if (dicAttr == null) {
                UnityEngine.Debug.LogWarning($"属性为空！");
                return;
            }

            List<Stats> list = new(dicAttr.Values);
            list.Sort((a, b) => (int)a.attr - (int)b.attr);

            string strInfo = titleInfo + "\n";
            for (int i = 0; i < list.Count; i++)
            {
                strInfo += $"{list[i].attr.ToString()}_{list[i].Num}\n";
            }
            UnityEngine.Debug.LogWarning(strInfo);
        }

        //转换成基础属性
        private void convertBaseStats(List<Stats> originalData, List<Stats> result, Dictionary<pb.Attribute, Stats> dicAttr)
        {
            Dictionary<pb.Attribute, int> temp = new();
            foreach (var item in originalData)
            {
                temp.Add(item.attr, item.Num);
            }

            //为了方便装备新增属性不在服务器下发的属性列表中，此处构造一个全部的属性列表
            foreach (pb.Attribute attr in Enum.GetValues(typeof(pb.Attribute)))
            {
                if (attr == pb.Attribute.INVALID_Attribute || attr == pb.Attribute.MAX_Attribute) continue;

                var stats = new Stats() { attr = attr, Num = 0 };
                result.Add(stats);

                //当前属性服务器已有数据，则直接赋值
                if (temp.ContainsKey(attr))
                {
                    stats.Num = temp[attr];
                }

                dicAttr.Add(attr, stats);
            }
        }

        //转换成基础属性
        private void convertToBaseStats(List<Stats> panelData, List<Stats> result, Dictionary<pb.Attribute, Stats> dicAttr)
        {
            Dictionary<pb.Attribute, int> temp = new();
            foreach (var item in panelData)
            {
                temp.Add(item.attr, item.Num);
            }

            //为了方便装备新增属性不在服务器下发的属性列表中，此处构造一个全部的属性列表
            foreach (pb.Attribute attr in Enum.GetValues(typeof(pb.Attribute)))
            {
                if (attr == pb.Attribute.INVALID_Attribute || attr == pb.Attribute.MAX_Attribute) continue;

                var stats = new Stats() { attr = attr, Num = 0 };
                result.Add(stats);

                //当前属性服务器已有数据，则直接赋值
                if (temp.ContainsKey(attr))
                {
                    if (AdjustAttr.Contains(attr))
                    {
                        var percent = AttrRelationShip[attr];
                        if (temp.ContainsKey(percent))
                        {
                            stats.Num = (int)(temp[attr] / (1 + temp[percent] / BaseRate));
                        }
                        else
                            stats.Num = temp[attr];
                    }
                    else
                        stats.Num = temp[attr];
                }

                dicAttr.Add(attr, stats);
            }
        }

        //遍历装备属性并移除当前装备所加的属性
        private void removeEquipAttrFromStatsDic(Dictionary<pb.Attribute, int> combinItemAttr, Dictionary<pb.Attribute, Stats> dicAttr, pb.ItemInfo itemInfo)
        {
            int quality = (int)itemInfo.Equip.Quality;
            float rate = quality > 0 ? CFG_XMLConfig.xmlConfig.EquipQualityDic[quality]._value : 1;
            //ConfigXmlDocument

            HashSet<pb.Attribute> opedAttr = new();
            foreach (var item in combinItemAttr)
            {
                if (!dicAttr.ContainsKey(item.Key)) continue; //在现有服务器的属性列表中没有找到相关属性

                if (opedAttr.Contains(item.Key)) continue; //已处理过

                if (AdjustAttr.Contains(item.Key))
                {
                    (int value, int percentage) = GetAdjustAttributeValue(dicAttr, item.Key);

                    if (value == 0) continue; //没有找到相关属性
                    double realPerc = percentage / BaseRate;
                    var baseVale = value / (1 + realPerc);
                    var newValue = (baseVale - item.Value)* (1 + realPerc);
                    dicAttr[item.Key].Num = (int)newValue;
                }
                else
                {
                    if (dicAttr[item.Key].Num - item.Value <= 0)
                        dicAttr[item.Key].Num = 0;
                    else
                        dicAttr[item.Key].Num -= item.Value;
                }
            }
        }

        //遍历装备属性并增加当前装备所加的属性
        private void addAttrToStatsDic(Dictionary<pb.Attribute, int> combinItemAttr, Dictionary<pb.Attribute, Stats> dicAttr)
        {
            foreach (var item in combinItemAttr)
            {
                if (!dicAttr.ContainsKey(item.Key)) continue; //没有找到属性

                if (AdjustAttr.Contains(item.Key))
                {
                    (int value, int percentage) = GetAdjustAttributeValue(dicAttr, item.Key);

                    if (value == 0) continue; //没有找到相关属性

                    double realPerc = percentage / BaseRate;
                    double baseValue = value / (1 + realPerc);
                    var newValue = (baseValue + item.Value) * (1 + realPerc);
                    dicAttr[item.Key].Num = (int)newValue;
                }
                else
                {
                    dicAttr[item.Key].Num += item.Value;
                }
            }
        }


        //获取装备属性
        private Dictionary<pb.Attribute, int> getItemAttributes(pb.ItemInfo itemInfo)
        {
            Dictionary<pb.Attribute, int> combinItemAttr = new();
            if (itemInfo == null) return combinItemAttr;
            var cfgInfo = CFG_Item.GetById(itemInfo.CfgId);

            int quality = (int)itemInfo.Equip.Quality;
            float rate = quality > 0 ? CFG_XMLConfig.xmlConfig.EquipQualityDic[quality]._value : 1;

            //获取基础属性
            foreach (var item in CFG_XMLConfig.xmlConfig.EquipDisplayDic)
            {
                var value = cfgInfo.GetType().GetField(item.Key).GetValue(cfgInfo);
                if (value == null) continue;

                int attrValue = (int)Math.Round((int)value * rate);
                if(attrValue <=0) continue;

                pb.Attribute attr = (pb.Attribute)CFG_XMLConfig.xmlConfig.EquipDisplayToSeverDic[item.Key];
                if (combinItemAttr.ContainsKey(attr))
                    combinItemAttr[attr] += attrValue;
                else
                    combinItemAttr.Add(attr, attrValue);
            }

            //获取高阶属性
            foreach (var item in itemInfo.Equip.Attributes)
            {
                if (combinItemAttr.ContainsKey(item.attr))
                    combinItemAttr[item.attr] += item.Num;
                else
                    combinItemAttr.Add(item.attr, item.Num);
            }

            return combinItemAttr;
        }

        private (string, string) GetItemBaseAttribute(int cfgID, int Quality)
        {
            var cfgInfo = CFG_Item.GetById(cfgID);

            float equipValue = Quality > 0 ? CFG_XMLConfig.xmlConfig.EquipQualityDic[Quality]._value : 1;
            string Str1 = "";
            string Str2 = "";

            foreach (var item in CFG_XMLConfig.xmlConfig.EquipDisplayDic)
            {
                if (cfgInfo.GetType().GetField(item.Key).GetValue(cfgInfo) != null && (int)cfgInfo.GetType().GetField(item.Key).GetValue(cfgInfo) > 0)
                {
                    int _value = (int)Math.Round((int)cfgInfo.GetType().GetField(item.Key).GetValue(cfgInfo) * equipValue);
                    //int qulite = CFG_RandomDistrib.GetColor(Quality > 0 ? Quality :(int)cfgInfo.Quality , CFG_XMLConfig.xmlConfig.EquipDisplayToSeverDic[item.Key], _value);
                    string arrtStr = CFG_XMLConfig.GetAttributeStrByID((int)CFG_XMLConfig.xmlConfig.EquipDisplayToSeverDic[item.Key], _value, true);
                }
            }
            return (Str1, Str2);
        }

        private (string, string) GetItemBestAttribute(int cfgID, EquipInfo eInfo)
        {
            var cfgInfo = CFG_Item.GetById(cfgID);
            string bestStr1 = "";
            string bestStr2 = "";
            int idx = 0;
            if (eInfo != null)
            {
                Dictionary<int, int> BestDic = new Dictionary<int, int>();
                foreach (var item in eInfo.Attributes)
                {
                    BestDic.Add((int)item.attr, item.Num);
                }

                foreach (var item in CFG_XMLConfig.xmlConfig.RolePanelList)
                {
                    if (BestDic.ContainsKey(item.SID) && BestDic[item.SID] > 0)
                    {
                        //int qulite = CFG_RandomDistrib.GetColor((int)eInfo.Quality, item.SID, BestDic[item.SID]);
                        string arrtStr = CFG_XMLConfig.GetAttributeStrByID(item.SID, BestDic[item.SID], true);
                        if (idx % 2 == 0)
                            bestStr1 = $"{bestStr1}{arrtStr}\n";
                        else
                            bestStr2 = $"{bestStr2}{arrtStr}\n";
                    }
                }
            }
            return (bestStr1, bestStr2);
        }

        
        // 获取属性值及属性对应的百分比值
        private (int curVal, int curPrecentage) GetAdjustAttributeValue(Dictionary<pb.Attribute, Stats> dicAttr, pb.Attribute attr)
        {
            if (!dicAttr.ContainsKey(attr)) return (0, 0);

            if (!AttrRelationShip.ContainsKey(attr)) return (0, 0);

            pb.Attribute percentage = AttrRelationShip[attr];
            if (!dicAttr.ContainsKey(percentage)) return (0, 0);

            return (dicAttr[attr].Num, dicAttr[percentage].Num);

        }
        #endregion

    }
}