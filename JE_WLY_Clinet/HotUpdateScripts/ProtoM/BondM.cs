using HotUpdateScripts.Config;
using JEngine.Singleton;
using MetaFramework.UI;
using System.Collections.Generic;

namespace HotUpdateScripts.ProtoM
{
    public enum BondState
    {
        None = 0,
        Lock = 1,
        Active = 2,
    }

    public class BondItem
    {
        public int bondID;
        public CFGItem_PartnerLink cfgData;
        public BondState bondState;

        public BondItem(CFGItem_PartnerLink cfg, BondState state)
        {
            RefreshData(cfg, state);
        }

        public void RefreshData(CFGItem_PartnerLink cfg, BondState state)
        {
            cfgData = cfg;
            bondID = cfg.BondID;
            bondState = state;
        }
    }

    public class BondM : SingletonTemplate<BondM>
    {
        private Dictionary<int, BondItem> partnerLinkDict = new Dictionary<int, BondItem>();

        private Dictionary<int, int> newBondDict = new Dictionary<int, int>();

        public int GetRedNum()
        {
            return newBondDict.Count;
        }


        private void LoadPartnerLinkDict()
        {
            partnerLinkDict.Clear();
            var allPartnerLinks = CFG_PartnerLink.GetAllPartnerLink();
            foreach (var pairs in allPartnerLinks)
            {
                pairs.Value.Sort((x, y) => { return x.ID.CompareTo(y.ID); });
                // 判断是否激活
                var firstPartnerLink = pairs.Value[0];
                var isActive = false;
                var minStar = 9999;
                for (int i = 0; i < firstPartnerLink.PartnerIDs.Length; i++)
                {
                    var partnerInfo = PartnerM.Instance.GetOnePartnerInfoByCFGID(firstPartnerLink.PartnerIDs[i]);
                    if (partnerInfo != null)
                    {
                        isActive = true;
                        if (partnerInfo.Star < minStar)
                            minStar = partnerInfo.Star;
                    }
                    else
                    {
                        isActive = false;
                        minStar = 0;
                        break;
                    }
                }

                if (!isActive)
                    partnerLinkDict.Add(firstPartnerLink.BondID, new BondItem(firstPartnerLink, BondState.Lock));
                else
                    partnerLinkDict.Add(firstPartnerLink.BondID, new BondItem(pairs.Value[minStar], BondState.Active));
            }

        }

        public Dictionary<int, BondItem> GetPartnerLinkDict()
        {
            LoadPartnerLinkDict();

            return partnerLinkDict;
        }

        public void CheckAddRed(int cfgId)
        {
            if (partnerLinkDict == null || partnerLinkDict.Count == 0)
                LoadPartnerLinkDict();

            var bonds = CFG_PartnerLink.GetPartnerLinkBond(cfgId);
            if (bonds == null || bonds.Count == 0)
                return;
            //UnityEngine.Debug.LogError($"------------------    {cfgId}, {bonds.Count}");

            for (int i = 0; i < bonds.Count; i++)
            {
                if (partnerLinkDict.TryGetValue(bonds[i], out var bondItem))
                {
                    //UnityEngine.Debug.LogError($"{bondItem.bondID}, {bondItem.bondState}");
                    //if (bondItem.bondState == BondState.Active)
                    //    continue;

                    var isActive = false;
                    var minStar = 9999;
                    for (int j = 0; j < bondItem.cfgData.PartnerIDs.Length; j++)
                    {
                        var partnerInfo = PartnerM.Instance.GetOnePartnerInfoByCFGID(bondItem.cfgData.PartnerIDs[j]);
                        if (partnerInfo != null)
                        {
                            //UnityEngine.Debug.LogError($"{partnerInfo.CfgId}, {partnerInfo.Star}, {minStar}");
                            isActive = true;
                            if (partnerInfo.Star < minStar)
                                minStar = partnerInfo.Star;
                        }
                        else
                        {
                            isActive = false;
                            minStar = 0;
                            break;
                        }
                    }

                    //UnityEngine.Debug.LogError($"{cfgId}, bondID:{bondItem.bondID}, isActive:{isActive},redNum:{redNum},bondState:{bondItem.bondState} ");
                    if (isActive)
                    {
                        if (bondItem.cfgData.Star != minStar || bondItem.bondState != BondState.Active)
                        {
                            bondItem.bondState = BondState.Active;
                            if (!newBondDict.ContainsKey(bondItem.bondID))
                                newBondDict.Add(bondItem.bondID, 1);

                            var ncfgId = bondItem.cfgData.ID + minStar - bondItem.cfgData.Star;
                            bondItem.cfgData = CFG_PartnerLink.GetPartnerLinkById(ncfgId);
                            //UnityEngine.Debug.LogError($"redNum++, ncfgId: {ncfgId}");
                        }
                    }
                }
            }

            if (newBondDict.Count > 0)
                RedPointManager.Instance.RPS.Dispatch(RedPointConst.Role_Bond.ToString(), newBondDict.Count);
        }

        public bool IsNew(int bondId)
        {
            if (newBondDict.ContainsKey(bondId))
                return true;

            return false;
        }

        public void RemoveNewRed(int bondId)
        {
            if (newBondDict.ContainsKey(bondId))
                newBondDict.Remove(bondId);

            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Role_Bond.ToString(), newBondDict.Count);
        }
    }
}