using JEngine.Helper;
using System.Collections.Generic;

public class CFGItem_RandomDistribV2 : CFGItemBase
{
    public int ClientMark;
    public override long GetIdx()
    {
        return ClientMark;
    }
    public int AttributeType;

    public int MaxLevel;

    public int RefingSegment;

    public string AttributeName;

    public int ValueMin;
    public int ValueMax;

    public int QualityColor;

}

/// <summary>
/// RandomDistrib 配置表
/// </summary>
public class CFG_RandomDistribV2 : CFGTableBase<CFGItem_RandomDistribV2>
{

    private static CFG_RandomDistribV2 _AllCFG;
    private static CFG_RandomDistribV2 AllCFG
    {
        get
        {
            _AllCFG ??= ExcelDataManager.Instance.GetExcelTable<CFG_RandomDistribV2, CFGItem_RandomDistribV2>();
            return _AllCFG;
        }
    }

    public static CFG_RandomDistribV2 Get()
    {
        return AllCFG;
    }

    public static CFGItem_RandomDistribV2 GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    private static Dictionary<int, List<CFGItem_RandomDistribV2>> _DicRandom;

    /// <summary>
    /// 根据属性类型（AttributeId）分类的属性列表
    /// </summary>
    public static Dictionary<int, List<CFGItem_RandomDistribV2>> DicRandom
    {
        get
        {
            if (_DicRandom == null)
            {
                _DicRandom = new Dictionary<int, List<CFGItem_RandomDistribV2>>();
                foreach (var item in Get().DicDatas.Values)
                {
                    if (!_DicRandom.ContainsKey(item.AttributeType))
                    {
                        _DicRandom.Add(item.AttributeType, new List<CFGItem_RandomDistribV2>());
                    }
                    _DicRandom[item.AttributeType].Add(item);
                }
                foreach (var item in _DicRandom)
                {
                    item.Value.Sort((a, b) => { return a.RefingSegment.CompareTo(b.RefingSegment); });
                }
            }
            return _DicRandom;
        }
    }

    public static string GetAttributeName(int attributeType)
    {
        if (DicRandom.ContainsKey(attributeType) && DicRandom[attributeType].Count > 0)
        {
            return DicRandom[attributeType][0].AttributeName;
        }
        return "";
    }

    /// <summary>
    /// 根据属性id和属性值获取对应的配置
    /// </summary>
    /// <param name="attributeId">属性id</param>
    /// <param name="value">属性值</param>
    /// <returns>cfg item</returns>
    public static CFGItem_RandomDistribV2 GetItem(int attributeId, int value)
    {
        if (DicRandom.ContainsKey(attributeId))
        {
            foreach (var item in DicRandom[attributeId])
            {
                if (value >= item.ValueMin && value <= item.ValueMax)
                {
                    return item;
                }
            }
        }
        return null;
    }

    /// <summary>
    /// 属性最大值
    /// </summary>
    /// <param name="attributeId">属性id</param>
    /// <returns>bool</returns>
    public static CFGItem_RandomDistribV2 GetMaxValue(int attributeId)
    {
        if (DicRandom.ContainsKey(attributeId))
        {
            int count = DicRandom[attributeId].Count;
            return DicRandom[attributeId][count - 1];
        }
        return null;
    }

    /// <summary>
    /// 属性值是否达到最大值
    /// </summary>
    /// <param name="attributeId">属性id</param>
    /// <param name="value">属性值</param>
    /// <returns>bool</returns>
    public static bool IsMaxValue(int attributeId, int value)
    {
        var item = GetMaxValue(attributeId);
        if (item != null && value >= item.ValueMax)
        {
            return true;
        }
        return false;
    }
}
