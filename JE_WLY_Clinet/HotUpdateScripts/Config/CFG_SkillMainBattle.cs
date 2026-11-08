using JEngine.Helper;
using System.Collections.Generic;
using UnityEngine;

public class CFGItem_SkillMainBattle : CFGItemBase
{
    public int ID;

    public override long GetIdx()
    {
        return ID;
    }

    public int Vocation; // 职业
    public string Name; // 技能名称
    public string Combo; // 技能组合
    public int Damage; // 技能伤害
    public int Level; // 角色等级限制
    public string SkillName; // 角色动作
    public int ActSpeed; // 角色动作速度
    public int CD; // 技能冷却（毫秒
    public int Effectnum; // 作用数量
    public int HitNum; // 伤害次数
    public int HitTime; // 伤害时长
    public string EffectName; // 技能效果
    public int EffectScaleInfo; // 技能缩放
    public int[] EffectPositionInfo; // 技能特效信息
    public int EffectStartDelay; // 起手延迟(毫秒
    public int MaxDistance; // 最大距离
    public string HitEffect; // 怪物受击技能
    public int HitEffectScaleInfo; // 受击缩放
    public int[] HitType; // 怪物受击特效类型
    public int[] HitDelay; // 怪物受击特效延时（毫秒）
    public int[] HitPositionInfo; // 怪物受击特效位置信息
    public int[] ColliderPositionInfo; // 碰撞盒子位置信息
    public int ColliderType; // 碰撞盒子类型信息
    public int[] ColliderSize; // 碰撞盒大小
    public int[] ColliderOffset; // 碰撞盒偏移
    public int ColliderDelay; // 碰撞盒出现延迟（毫秒）
    public string Image; // 技能名图
}

public class CFG_SkillMainBattle : CFGTableBase<CFGItem_SkillMainBattle>
{
    private static List<CFGItem_SkillMainBattle> norSpecialSkills = new List<CFGItem_SkillMainBattle>();
    private static List<CFGItem_SkillMainBattle> bigSpecialSkills = new List<CFGItem_SkillMainBattle>();

    private static CFG_SkillMainBattle _AllCFG;
    private static CFG_SkillMainBattle AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_SkillMainBattle, CFGItem_SkillMainBattle>();
            }
            return _AllCFG;
        }
    }

    public static CFG_SkillMainBattle Get()
    {
        return AllCFG;
    }

    public static CFGItem_SkillMainBattle GetById(int id)
    {
        if (!AllCFG.DicDatas.ContainsKey(id))
            return null;
        return AllCFG.DicDatas[id];
    }

    public static List<CFGItem_SkillMainBattle> GetNormalSkillsByRoleVocation(int vocation, int lv)
    {
        norSpecialSkills.Clear();
        bigSpecialSkills.Clear();
        var temp = new List<CFGItem_SkillMainBattle>();
        var list = Get().DicList;
        for (int i = 0; i < list.Count; i++)
        {
            if (list[i].Vocation == vocation && string.IsNullOrEmpty(list[i].Combo) && list[i].Level <= lv)
                temp.Add(list[i]);

            if (list[i].Vocation == vocation && !string.IsNullOrEmpty(list[i].Combo) && list[i].Level <= lv)
            {
                norSpecialSkills.Add(list[i]);
            }

            if (list[i].Vocation == 5 && !string.IsNullOrEmpty(list[i].Combo) && list[i].Level <= lv)
            {
                bigSpecialSkills.Add(list[i]);
            }    
        }
        return temp;
    }

    public static CFGItem_SkillMainBattle GetOneNorSpecialSkill()
    {
        if (norSpecialSkills == null || norSpecialSkills.Count == 0)
            return null;

        if(norSpecialSkills.Count == 1)
            return norSpecialSkills[0];

        var index = Random.Range(0, norSpecialSkills.Count);
        return norSpecialSkills[index];
    }

    public static CFGItem_SkillMainBattle GetOneBigSpecialSkill()
    {
        if (bigSpecialSkills == null || bigSpecialSkills.Count == 0)
            return null;

        if (bigSpecialSkills.Count == 1)
            return bigSpecialSkills[0];

        var index = Random.Range(0, bigSpecialSkills.Count);
        return bigSpecialSkills[index];
    }
}