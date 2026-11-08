using JEngine.Helper;
using System.Collections.Generic;
using System.Linq;

namespace HotUpdateScripts.Config
{
    public class CFGItem_PartnerLink : CFGItemBase
    {
        public int ID;

        public override long GetIdx()
        {
            return ID;
        }

        public int BondID;

        public string BondName;

        public int[] PartnerIDs; // 伙伴列表

        public int Star; // 最低星级

        public Dictionary<string, int>[] Attribute; // 属性
    }

    public class CFG_PartnerLink : CFGTableBase<CFGItem_PartnerLink>
    {
        private static CFG_PartnerLink _AllCFG;
        private static Dictionary<int, List<CFGItem_PartnerLink>> partnerLinkDict = new Dictionary<int, List<CFGItem_PartnerLink>>();
        private static Dictionary<int, HashSet<int>> partnerBondDict = new Dictionary<int, HashSet<int>>(); // key 伙伴ID, value为bondid

        private static CFG_PartnerLink AllCFG
        {
            get
            {
                if (_AllCFG == null)
                    _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_PartnerLink, CFGItem_PartnerLink>();

                return _AllCFG;
            }

        }

        public static CFG_PartnerLink Get()
        {
            return AllCFG;
        }

        public static Dictionary<int, List<CFGItem_PartnerLink>> GetAllPartnerLink()
        {
            if (partnerLinkDict == null || partnerLinkDict.Count == 0)
            {
                partnerLinkDict.Clear();
                partnerBondDict.Clear();

                var allDatas = AllCFG.DicList;
                for (int i = 0; i < allDatas.Count; i++)
                {
                    var bondId = allDatas[i].BondID;
                    if (!partnerLinkDict.ContainsKey(bondId))
                        partnerLinkDict.Add(bondId, new List<CFGItem_PartnerLink>());

                    partnerLinkDict[bondId].Add(allDatas[i]);

                    for (int j = 0; j < allDatas[i].PartnerIDs.Length; j++)
                    {
                        var id = allDatas[i].PartnerIDs[j];
                        if (!partnerBondDict.ContainsKey(id))
                            partnerBondDict.Add(id, new HashSet<int>());

                        partnerBondDict[id].Add(bondId);
                    }
                }
            }

            return partnerLinkDict;
        }

        public static List<CFGItem_PartnerLink> GetPartnerLinkByBondId(int bondId)
        {
            if (partnerLinkDict == null || partnerLinkDict.Count == 0)
                GetAllPartnerLink();

            if (partnerLinkDict.TryGetValue(bondId, out var list))
                return list;

            return null;
        }

        public static CFGItem_PartnerLink GetPartnerLinkById(int cfgId)
        {
            if (AllCFG.DicDatas.ContainsKey(cfgId))
                return AllCFG.DicDatas[cfgId];

            return null;
        }

        public static List<int> GetPartnerLinkBond(int cfgId)
        {
            if (partnerBondDict == null || partnerBondDict.Count == 0)
                GetAllPartnerLink();

            if (partnerBondDict.TryGetValue(cfgId, out var list))
                return list.ToList();

            return null;
        }
    }
}