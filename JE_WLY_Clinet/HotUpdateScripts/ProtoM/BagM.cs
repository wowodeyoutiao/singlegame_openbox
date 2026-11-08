using JEngine.Helper;
using JEngine.Singleton;
using pb;

using System.Collections.Generic;
using System.Linq;

namespace HotUpdateScripts.ProtoM
{
    /// <summary>
    /// 包裹信息
    /// </summary>
    public class BagM : SingletonTemplate<BagM>
    {
        private List<ItemInfo> _itemList = new List<ItemInfo>();
        //总道具
        public List<ItemInfo> ItemList
        {
            get
            {
                _itemList.Clear();
                var list = UidItemsDic.Values.ToArray();
                for (int i = 0; i < list.Length; i++)
                {
                    if (!CheckItemTime(list[i]) && CFG_Item.GetById(list[i].CfgId).Type != (int)CFG_Item.ItemType.NoBag)
                    {
                        _itemList.Add(list[i]);
                    }
                }
                _itemList.Sort(ItemCompare);
                return _itemList;
            }
        }

        /// <summary>
        /// 检测道具是否到期 false 为未到期
        /// </summary>
        /// <param name="_item"></param>
        /// <returns></returns>
        private bool CheckItemTime(ItemInfo _item)
        {
            if (_item.CreateTime == 0 || CFG_Item.GetById(_item.CfgId).Timeliness == null)
            {
                return false;
            }
            if (CFG_Item.GetById(_item.CfgId).Timeliness.type == (int)CFG_Item.ItemTimeType.EndTimeDay)
            {
                long endTime = _item.CreateTime + 86400 * CFG_Item.GetById(_item.CfgId).Timeliness.time;
                return GameUtil.Instance.GetCurrentTimeSeconds() >= endTime;
            }
            else
            {
                return GameUtil.Instance.GetCurrentTimeSeconds() >= CFG_Item.GetById(_item.CfgId).Timeliness.endtime;
            }
        }

        private List<ItemInfo> _equipItemList = new List<ItemInfo>();
        private bool _isEquipListDirty = false;

        public bool GetEquipChange()
        {
            return _isEquipListDirty;
        }
        //装备
        public List<ItemInfo> GetEquipItemList(bool isSolt = false)
        {
            if (_isEquipListDirty)
            {
                _isEquipListDirty = false;

                _equipItemList.Clear();
                var list = UidItemsDic.Values.ToArray();
                for (int i = 0; i < list.Length; i++)
                {
                    if (CFG_Item.GetById(list[i].CfgId).Type == (int)CFG_Item.ItemType.Equip && !CheckItemTime(list[i]) && list[i].Num > 0)
                    {
                        _equipItemList.Add(list[i]);
                    }
                }
                //if (isSolt)
                _equipItemList.Sort(ItemCompare);
            }

            return _equipItemList;
        }

        private List<ItemInfo> _otherItemList = new List<ItemInfo>();
        public List<ItemInfo> GetOtherItemList(bool isSolt = false)
        {
            _otherItemList.Clear();
            var list = UidItemsDic.Values.ToArray();
            for (int i = 0; i < list.Length; i++)
            {
                if (CFG_Item.GetById(list[i].CfgId).Type != (int)CFG_Item.ItemType.NoBag && CFG_Item.GetById(list[i].CfgId).Type != (int)CFG_Item.ItemType.Equip
                    && CFG_Item.GetById(list[i].CfgId).Type != (int)CFG_Item.ItemType.Jewel && CFG_Item.GetById(list[i].CfgId).Type != (int)CFG_Item.ItemType.PartnerPice
                    && !CheckItemTime(list[i]))
                {
                    _otherItemList.Add(list[i]);
                }
            }
            if (isSolt)
                _otherItemList.Sort(ItemCompare);
            return _otherItemList;
        }

        private List<ItemInfo> _JewelItemList = new List<ItemInfo>();
        /// <summary>
        /// 宝石
        /// </summary>
        public List<ItemInfo> JewelItemListList(bool isSolt = false)
        {
            _JewelItemList.Clear();
            var list = UidItemsDic.Values.ToArray();
            for (int i = 0; i < list.Length; i++)
            {
                if (CFG_Item.GetById(list[i].CfgId).Type == (int)CFG_Item.ItemType.Jewel && !CheckItemTime(list[i]))
                {
                    _JewelItemList.Add(list[i]);
                }
            }
            if (isSolt)
                _JewelItemList.Sort(ItemCompare);
            return _JewelItemList;
        }


        private List<ItemInfo> _PartnerPicesList = new List<ItemInfo>();
        /// <summary>
        /// 伙伴碎片
        /// </summary>
        public List<ItemInfo> GetPartnerPicesList(bool isSolt = false)
        {
            _PartnerPicesList.Clear();
            var list = UidItemsDic.Values.ToArray();
            for (int i = 0; i < list.Length; i++)
            {
                if (CFG_Item.GetById(list[i].CfgId).Type == (int)CFG_Item.ItemType.PartnerPice && !CheckItemTime(list[i]))
                {
                    _PartnerPicesList.Add(list[i]);
                }
            }
            if (isSolt)
                _PartnerPicesList.Sort(ItemCompare);
            return _PartnerPicesList;
        }

        /// <summary>
        /// UID 对应道具列表
        /// </summary>
        public Dictionary<long, ItemInfo> UidItemsDic = new Dictionary<long, ItemInfo>();

        /// <summary>
        /// cfg 对应的所有UID信息
        /// </summary>
        public Dictionary<int, HashSet<long>> CFGtoUidDic = new Dictionary<int, HashSet<long>>();

        public HashSet<long> TitileNewHashSet = new HashSet<long>();

        /// <summary>
        /// cfg 对应的包裹内道具总数列表
        /// </summary>
        public Dictionary<int, long> CIDAndNumDic = new Dictionary<int, long>();

        public List<DailyItemUsed> dailyUsedItems = new List<DailyItemUsed>();

        public void SetBagInfo(List<ItemInfo> bagsInfo)
        {
            _isEquipListDirty = true;
            TitileNewHashSet.Clear();
            UidItemsDic.Clear();
            CIDAndNumDic.Clear();
            CFGtoUidDic.Clear();

            for (int i = 0; i < bagsInfo.Count; i++)
            {
                //标记道具 伙伴整卡 宠物 不进入包裹
                if (CFG_Item.GetById(bagsInfo[i].CfgId).Type == 0 || CFG_Item.GetById(bagsInfo[i].CfgId).Type == 8 || CFG_Item.GetById(bagsInfo[i].CfgId).Type == 9)
                {
                    continue;
                }

                if (UidItemsDic.ContainsKey(bagsInfo[i].Uid))
                {
                    UidItemsDic[bagsInfo[i].Uid] = bagsInfo[i];
                }
                else
                {
                    UidItemsDic.Add(bagsInfo[i].Uid, bagsInfo[i]);
                }

                if (!CFGtoUidDic.ContainsKey(bagsInfo[i].CfgId))
                {
                    CFGtoUidDic.Add(bagsInfo[i].CfgId, new HashSet<long>());
                }
                CFGtoUidDic[bagsInfo[i].CfgId].Add(bagsInfo[i].Uid);

                AddOneCIDAndNumDic(bagsInfo[i].CfgId, bagsInfo[i].Num);
            }
        }

        public void SetDailyUsedItems(List<DailyItemUsed> dailyItems)
        {
            dailyUsedItems.Clear();

            if (dailyItems != null)
            {
                dailyUsedItems = dailyItems;
            }
        }

        /// <summary>
        /// 道具每日已使用次数
        /// </summary>
        /// <param name="cfgId">道具id</param>
        /// <returns></returns>
        public long GetDailyUsedItem(int cfgId)
        {
            var item = dailyUsedItems.Find(item => item.CfgId == cfgId);
            if (item == null) return 0;
            return item.Num;
        }

        private void AddOneCIDAndNumDic(int cfgId, long Num)
        {
            if (!CIDAndNumDic.ContainsKey(cfgId))
            {
                CIDAndNumDic.Add(cfgId, 0);
            }
            CIDAndNumDic[cfgId] += Num;

            if (CFG_Item.GetById(cfgId).Type == (int)CFG_Item.ItemType.Equip)
            {
                _isEquipListDirty = true;
            }
        }

        /// <summary>
        /// 更具CFG ID 获取数量
        /// </summary>
        /// <param name="cfgID"></param>
        /// <returns></returns>
        public long GetNumByCfgID(int cfgID)
        {
            if (cfgID == (int)CFG_Item.ItemID.ZhanHun)
                return RoleM.Instance.Soul;
            else if (cfgID == (int)CFG_Item.ItemID.Gold)
                return RoleM.Instance.Coin;
            else if (cfgID == (int)CFG_Item.ItemID.DianQuan)
                return RoleM.Instance.Coupon;
            else if (cfgID == (int)CFG_Item.ItemID.SLPZ)
                return RoleM.Instance.TrialVoucher;
            else if (cfgID == (int)CFG_Item.ItemID.PetEXP)
                return RoleM.Instance.PetUsedExp;
            else if (cfgID == (int)CFG_Item.ItemID.JSJZ)
                return RoleM.Instance.BoxQuick;
            else if (cfgID == (int)CFG_Item.ItemID.EXP)
                return RoleM.Instance.Exp;
            else if (cfgID == (int)CFG_Item.ItemID.GuildContribution)
                return RoleM.Instance.GuildContribution;
            else if (cfgID == (int)CFG_Item.ItemID.GuildWish)
                return RoleM.Instance.GuildWish;
            else if (cfgID == (int)CFG_Item.ItemID.StarWishCoin)
                return RoleM.Instance.StarWishCoin;

            if (CIDAndNumDic.ContainsKey(cfgID))
                return CIDAndNumDic[cfgID];

            return 0;
        }


        /// <summary>
        /// 添加资源
        /// </summary>
        /// <param name="bagsInfo"></param>
        public void AddItemsChange(List<ItemInfo> bagsInfo)
        {
            for (int i = 0; i < bagsInfo.Count; i++)
            {
                AddOneItemChange(bagsInfo[i]);

                if (!TitileNewHashSet.Contains(bagsInfo[i].CfgId) && CFG_Item.GetById(bagsInfo[i].CfgId).Type == (int)CFG_Item.ItemType.TitlePice)
                {
                    TitileNewHashSet.Add(bagsInfo[i].CfgId);
                }
            }
        }

        public bool IsGetTitlePiceNew(int cfgId)
        {
            return TitileNewHashSet.Contains(cfgId);
        }

        public bool RemoveTitlePiceNew(int cfgId)
        {
            return TitileNewHashSet.Remove(cfgId);
        }

        public void AddOneItemChange(ItemInfo bagInfo)
        {
            if (CFG_Item.GetById(bagInfo.CfgId).Type == 0 || CFG_Item.GetById(bagInfo.CfgId).Type == 8 || CFG_Item.GetById(bagInfo.CfgId).Type == 9)
            {
                return;
            }

            if (UidItemsDic.ContainsKey(bagInfo.Uid))
            {
                UidItemsDic[bagInfo.Uid].Num += bagInfo.Num;
            }
            else
            {
                UidItemsDic[bagInfo.Uid] = bagInfo;

                if (!CFGtoUidDic.ContainsKey(bagInfo.CfgId))
                {
                    CFGtoUidDic.Add(bagInfo.CfgId, new HashSet<long>());
                }

                CFGtoUidDic[bagInfo.CfgId].Add(bagInfo.Uid);
            }
            AddOneCIDAndNumDic(bagInfo.CfgId, bagInfo.Num);


        }

        /// <summary>
        /// 根据CFG 获取包裹内道具所有信息
        /// </summary>
        /// <param name="id"></param>
        /// <returns></returns>
        public ItemInfo[] GetItemsByCFGID(int id)
        {
            if (CFGtoUidDic.ContainsKey(id))
            {
                var uidInfo = CFGtoUidDic[id].ToArray();
                var items = new ItemInfo[uidInfo.Length];
                for (int i = 0; i < uidInfo.Length; i++)
                {
                    items[i] = UidItemsDic[uidInfo[i]];
                }
                return items;
            }
            return null;
        }

        public ItemInfo GetItemByCFGID(int id)
        {
            var valueList = GetItemsByCFGID(id);
            if (valueList != null && valueList.Length > 0)
            {
                return valueList[0];
            }
            return null;
        }

        /// <summary>
        /// 减少资源
        /// </summary>
        /// <param name="bagsInfo"></param>
        public void DelItemsChange(List<ItemInfo> bagsInfo)
        {
            for (int i = 0; i < bagsInfo.Count; i++)
            {
                if (UidItemsDic.ContainsKey(bagsInfo[i].Uid))
                {
                    UidItemsDic[bagsInfo[i].Uid].Num -= bagsInfo[i].Num;
                    if (UidItemsDic[bagsInfo[i].Uid].Num <= 0)
                    {
                        UidItemsDic.Remove(bagsInfo[i].Uid);

                        if (CFGtoUidDic.ContainsKey(bagsInfo[i].CfgId))
                            CFGtoUidDic[bagsInfo[i].CfgId].Remove(bagsInfo[i].Uid);
                    }
                    AddOneCIDAndNumDic(bagsInfo[i].CfgId, -bagsInfo[i].Num);
                }
            }
        }

        /// <summary>
        /// 包裹中的排序
        /// </summary>
        /// <param name="x"></param>
        /// <param name="y"></param>
        /// <returns></returns>
        public int ItemCompare(ItemInfo x, ItemInfo y)
        {
            //if (ReferenceEquals(x, y)) return 0;
            //if (x is null) return -1;
            //if (y is null) return 1;

            var cfgX = CFG_Item.GetById(x.CfgId);
            var cfgY = CFG_Item.GetById(y.CfgId);

            int r;

            // 1) Priority 降序
            r = cfgY.Priority.CompareTo(cfgX.Priority);
            if (r != 0) return r;

            // 2) CfgId 降序
            r = y.CfgId.CompareTo(x.CfgId);
            if (r != 0) return r;

            // 3) 都有 Equip 时按 Quality 降序
            if (cfgX.Type == (int)CFG_Item.ItemType.Equip && cfgY.Type == (int)CFG_Item.ItemType.Equip)
            {
                r = y.Equip.Quality.CompareTo(x.Equip.Quality);
                if (r != 0) return r;
            }

            // 4) Uid 降序
            return y.Uid.CompareTo(x.Uid);
        }


        //合成装备相同品质的快速填入
        public void GetAutoSelectEquipSyns(int allCount, ref List<UISynthesis.SelectSyn> selectSyns)
        {
            int id = selectSyns[0].id;
            long uid = selectSyns[0].Uid;
            selectSyns.Clear();

            void OneItemAdd(int _id, long _uid, ref List<UISynthesis.SelectSyn> selectSyns)
            {
                int AllNum = (int)GetNumByCfgID(_id);

                int num1 = AllNum / allCount;
                int num2 = AllNum % allCount;
                for (int i = 1; i <= allCount; i++)
                {
                    int _num = num1;
                    if (num2 >= i)
                    {
                        _num += 1;
                    }
                    if (_num > 0)
                    {
                        selectSyns.Add(new UISynthesis.SelectSyn { id = _id, num = _num, Uid = _uid });
                    }
                }
            }
            if (CFG_Item.GetById(id).Maxnum > 1 && GetNumByCfgID(id) >= allCount)
            {

                OneItemAdd(id, uid, ref selectSyns);
            }
            else
            {
                var raws = CFG_Compose.GetByRawId(id).Raw;
                var _ItemInfo = new List<ItemInfo>();

                for (int i = 0; i < raws.Length; i++)
                {
                    var items = GetItemsByCFGID(raws[i]);
                    if (items != null)
                    {
                        for (int j = 0; j < items.Length; j++)
                        {
                            _ItemInfo.Add(items[j]);
                        }
                    }
                }

                int _addCount = allCount;
                for (int i = 0; i < _ItemInfo.Count; i++)
                {
                    if (CFG_Item.GetById(_ItemInfo[i].CfgId).Maxnum > 1 && GetNumByCfgID(_ItemInfo[i].CfgId) >= allCount && i == 0)
                    {
                        OneItemAdd(_ItemInfo[i].CfgId, _ItemInfo[i].Uid, ref selectSyns);
                        break;
                    }
                    else
                    {
                        for (int j = 0; j < _ItemInfo[i].Num; j++)
                        {
                            _addCount--;
                            selectSyns.Add(new UISynthesis.SelectSyn { id = _ItemInfo[i].CfgId, num = 1, Uid = _ItemInfo[i].Uid });
                            if (_addCount == 0)
                                return;
                        }
                    }
                }
            }
        }

        //获取材料合成的所有信息
        public List<ItemInfo> GetSynListByType(int _type, List<UISynthesis.SelectSyn> selectSyns, bool solt = false)
        {
            var itemList = new List<ItemInfo>();

            if (selectSyns != null && selectSyns.Count != 0)
            {
                var raws = CFG_Compose.GetByRawId(selectSyns[0].id).Raw;

                HashSet<long> SelectInfo = new HashSet<long>();
                for (int i = 0; i < selectSyns.Count; i++)
                {
                    SelectInfo.Add(selectSyns[i].Uid);
                }

                for (int i = 0; i < raws.Length; i++)
                {
                    bool isHave = false;
                    int num = (int)GetNumByCfgID(raws[i]);
                    if (CFG_Item.GetById(raws[i]).Maxnum > 1)
                    {
                        for (int j = 0; j < selectSyns.Count; j++)
                        {
                            if (selectSyns[j].id == raws[i])
                            {
                                isHave = true;
                                num -= selectSyns[j].num;
                            }
                        }
                        if (!isHave || num > 0)
                        {
                            var item = GetItemByCFGID(raws[i]);
                            if (item != null)
                            {
                                item.Num = num;
                                itemList.Add(item);
                            }
                        }
                    }
                    else
                    {
                        var bagitems = GetItemsByCFGID(raws[i]);
                        if (bagitems != null)
                        {
                            for (int j = 0; j < bagitems.Length; j++)
                            {
                                if (!SelectInfo.Contains(bagitems[j].Uid))
                                {
                                    itemList.Add(bagitems[j]);
                                }
                            }
                        }
                    }
                }
            }
            else
            {
                List<ItemInfo> allList = null;
                if (_type == 3)
                    allList = JewelItemListList();
                else
                    allList = GetOtherItemList();

                Dictionary<int, int> CFG_IdxDic = new Dictionary<int, int>();
                int Idx = itemList.Count;

                for (int i = 0; i < allList.Count; i++)
                {
                    if (CFG_Item.GetById(allList[i].CfgId).SonType == _type && CFG_Item.GetById(allList[i].CfgId).ComposeMark == 1 && allList[i].Num > 0)
                    {
                        if (_type == 2)
                        {
                            itemList.Add(allList[i]);
                        }
                        else
                        {
                            if (!CFG_IdxDic.ContainsKey(allList[i].CfgId))
                            {
                                CFG_IdxDic.Add(allList[i].CfgId, Idx);
                                Idx++;
                                itemList.Add(allList[i]);
                            }
                            itemList[CFG_IdxDic[allList[i].CfgId]].Num = GetNumByCfgID(allList[i].CfgId);
                        }
                    }
                }
            }
            if (solt)
                itemList.Sort(ItemCompare);
            return itemList;
        }

        //自动分解设置
        public AutoDisassemblyInfo AutoDelSetInfo;

        public bool IsEquipFull()
        {
            return GetEquipItemList().Count >= CFG_XMLConfig.xmlConfig.BaseValue.bagNum.equipnum;
        }

        /// <summary>
        /// 检查装备空位是否够用
        /// </summary>
        /// <param name="needVacancy">需要的空余位置</param>
        /// <returns></returns>
        public bool CheckEquipVacancy(int needVacancy)
        {
            return needVacancy <= CFG_XMLConfig.xmlConfig.BaseValue.bagNum.equipnum - GetEquipItemList().Count;
        }

        /// <summary>
        /// 检查道具空位是否够用
        /// </summary>
        /// <param name="needVacancy">需要的空余位置</param>
        /// <returns></returns>
        public bool CheckItemVacancy(int needVacancy)
        {
            return needVacancy <= CFG_XMLConfig.xmlConfig.BaseValue.bagNum.itemnum - GetOtherItemList().Count;
        }
    }

}

