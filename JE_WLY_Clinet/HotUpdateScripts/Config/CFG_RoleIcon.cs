using JEngine.Helper;

using pb;

using System.Collections.Generic;
using System.Text;

public class CFGItem_RoleIcon : CFGItemBase
{
    public int ID;
    public string Name;
    public int ItemType;
    public Attribute[] Value;
    public int Period;
    public string Get;
    public string Path;
    public string Icon_ZS;
    public string Icon_FS;
    public string Icon_JS;
    public int Quality;

    public struct Attribute
    {
        public int key;
        public int value;
    }

    public override long GetIdx()
    {
        return ID;
    }

    public string GetAttributeString()
    {
        var sb = new StringBuilder();
        foreach (var attr in Value)
        {
            sb.Append(CFG_XMLConfig.GetAttributeStrByID(attr.key, attr.value));
            sb.Append("  ");
        }
        return sb.ToString();
    }

    public string GetIcon(int job)
    {
        string icon = job switch
        {
            1 => Icon_ZS,
            2 => Icon_JS,
            3 => Icon_FS,
            _ => Icon_ZS, //兜底：未知/0 职业也用战士图标，避免头像空白
        };

        return icon;
    }

}

/// <summary>
/// CFG_Name 配置表
/// </summary>
public class CFG_RoleIcon : CFGTableBase<CFGItem_RoleIcon>
{
    private static CFG_RoleIcon _AllCFG;
    private static CFG_RoleIcon AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_RoleIcon, CFGItem_RoleIcon>();
            }
            return _AllCFG;
        }
    }

    public static CFG_RoleIcon Get()
    {
        return AllCFG;
    }

    public static CFGItem_RoleIcon GetById(int id)
    {
        //单机服：新号 current_avatar_id / current_frame_id 可能为 0，直接索引会抛 KeyNotFoundException
        if (AllCFG.DicDatas.TryGetValue(id, out var item))
            return item;
        return null;
    }

    public static List<CFGItem_RoleIcon> GetListByType(AvatarType avatarType)
    {
        return GetListByType((int)avatarType);
    }

    public static List<CFGItem_RoleIcon> GetListByType(int type)
    {
        return AllCFG.DicList.FindAll(x => x.ItemType == type);
    }

    public static string GetAttributeString(int id)
    {
        var item = GetById(id);
        if (item == null) return string.Empty;

        var sb = new StringBuilder();
        foreach (var attr in item.Value)
        {
            sb.Append(CFG_XMLConfig.GetAttributeStrByID(attr.key, attr.value));
        }
        return sb.ToString();
    }

    public static (string atlas, string icon) GetIcon(int id, int job)
    {
        var cfg = GetById(id);
        if (cfg == null) return (string.Empty, string.Empty);

        string icon = job switch
        {
            1 => cfg.Icon_ZS,
            2 => cfg.Icon_JS,
            3 => cfg.Icon_FS,
            _ => string.Empty,
        };

        return (cfg.Path, icon);
    }

}
