using HotUpdateScripts;
using HotUpdateScripts.ProtoM;

using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using pb;

using Spine.Unity;

using System.Collections.Generic;

using UnityEngine;
using UnityEngine.UI;

public class UIOtherRole : UIBase<UIModelOtherRole, UIViewOtherRole>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIOtherRole>();
        });

        for (int i = 0; i < GetView().teamItems.Length; i++)
        {
            int idx = i;
            GetView().teamItems[i]._btn.onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                if (idx > 0 && CFG_XMLConfig.GetOpenSoltLv(idx) > observePlayerInfo.Level)
                {
                    return;
                }

                OnTeamItemClick(idx);
            });
        }

        GetView().btn_Attribute.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            // 属性展开、收起
            SetArrtInfo(!AttrShow);
        });

        for (int i = 0; i < GetView().root_Equips.childCount; i++)
        {
            int idx = i + 1;
            JUIExtensions.btn(GetView().root_Equips.Find("equip_" + idx)).onClick.AddListener(() =>
            {
                var itemInfo = TeamIdx == 0 ? observePlayerInfo.SlotInfo[idx].Item : GetPartnerSlotInfo(TeamIdx).Equips[idx].Item;
                if (itemInfo != null)
                {
                    var equipRefineData = GetOneEquipSlotData((CharacterSlot)(TeamIdx + 1), (EquipSlot)idx);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params
                    {
                        cfgId = itemInfo.CfgId,
                        equipInfo = itemInfo.Equip,
                        teamIdx = 0,
                        isIgnoreJob = true,
                        equipRefineData = equipRefineData,
                        equipSlotPos = idx,
                    });
                }
            });
        }
    }

    // 玩家信息
    ObservePlayerInfo observePlayerInfo;

    public override void OnEnter(params object[] args)
    {
        TeamIdx = 0;
        observePlayerInfo = args[0] as ObservePlayerInfo;

        SetPartnerLock();
        SetBaseRole();
    }

    private ObservePlayerInfo.PartnerStats GetPartnerState(long partnerUid)
    {
        return observePlayerInfo.PartnersStats.Find(v => v.PartnerUid == partnerUid);
    }

    public EquipRefineData GetOneEquipSlotData(CharacterSlot characterSlot, EquipSlot equipSlot)
    {
        CharacterRefineData refineData = observePlayerInfo.character_refines.Find(x => x.character_slot == characterSlot);
        if (refineData != null)
        {
            return refineData.equip_refines.Find(x => x.equip_slot == equipSlot);
        }
        return null;
    }

    private void SetPartnerLock()
    {
        for (int soltIdx = 1; soltIdx < GetView().teamItems.Length; soltIdx++)
        {
            bool isLock = CFG_XMLConfig.GetOpenSoltLv(soltIdx) > observePlayerInfo.Level;
            GetView().teamItems[soltIdx]._obj_lock.SetActive(isLock);
            if (isLock)
                GetView().teamItems[soltIdx]._txt_Lock.text = $"{CFG_XMLConfig.GetOpenSoltLv(soltIdx)}级解锁";
        }
    }

    //设置基础信息
    private void SetBaseRole()
    {
        for (int i = 0; i < GetView().teamItems.Length; i++)
        {
            if (i == 0)
            {
                GameUtil.Instance.SetJobIcon(observePlayerInfo.Job, GetView().teamItems[i]._img_job);
                //ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.RoleHeadAtlas, "img_" + observePlayerInfo.Job + observePlayerInfo.Sex, GetView().teamItems[i]._img_head);
                (int headId, int frameId) = GameUtil.Instance.GetHeadAndFrameId(observePlayerInfo.Avatar, observePlayerInfo.Frame);
                GameUtil.Instance.SetRoleHead(headId, RoleM.Instance.Job, GetView().teamItems[i]._img_head, GetView().teamItems[i]._img_head_quality);
                GameUtil.Instance.SetRoleHeadFrame(frameId, RoleM.Instance.Job, GetView().teamItems[i].img_frame);
                GetView().teamItems[i]._txt_lv.text = observePlayerInfo.Level.ToString();
            }
            else
            {
                int quality = 0;
                bool isShow = GetPartnerSlotInfo(i).PartnerUid != 0;
                GetView().teamItems[i]._img_head.gameObject.SetActive(isShow);
                GetView().teamItems[i]._img_job.gameObject.SetActive(isShow);
                GetView().teamItems[i]._obj_Lv.gameObject.SetActive(isShow);
                GetView().teamItems[i].txt_Star.gameObject.SetActive(isShow);
                if (isShow)
                {
                    var solt = GetPartnerSlotInfo(i);
                    var partner = GetPartnerInfo(solt.PartnerUid);
                    quality = CFG_Partner.GetById(partner.CfgId).Quality;
                    GetView().teamItems[i]._txt_lv.text = partner.Level.ToString();
                    GetView().teamItems[i].txt_Star.text = partner.Star.ToString();
                    GameUtil.Instance.SetJobIcon(CFG_Partner.GetById(partner.CfgId).job % 10, GetView().teamItems[i]._img_job);
                    GameUtil.Instance.SetItemIcon(CFG_Partner.GetById(partner.CfgId).ItemID, GetView().teamItems[i]._img_head);
                }
                GameUtil.Instance.SetPartnetQualityFrame(quality, GetView().teamItems[i].img_frame);
            }
        }

        GetView().txt_allPower.text = GameUtil.Instance.FormatPowerWithUnit(observePlayerInfo.TotalCombat);
        OnTeamItemClick(TeamIdx);
    }

    private void OnPetShow()
    {
        PetSlotInfo myPet = observePlayerInfo.Pets;
        GetView().img_Pet.Find("noPet").gameObject.SetActive(myPet.CfgId == 0);
        GetView().img_Pet.Find("Pet").gameObject.SetActive(myPet.CfgId != 0);

        if (myPet.CfgId != 0)
        {
            var root = GetView().img_Pet.Find("Pet");
            string pooName = "Pet_" + myPet.CfgId;
            _ = GameHelper.Instance.SetAsyncObj(pooName, GameConst.GetMonsterModelByID(pooName), (_mainRole) =>
            {
                if (root.childCount > 0)
                {
                    PoolHelper.Instance.RecycObj(root.GetChild(0).gameObject);
                }
                _mainRole.transform.SetParent(root);
                _mainRole.transform.localScale = Vector3.one;
                _mainRole.transform.localPosition = Vector3.zero;
                _mainRole.SetActive(true);
            });
        }

    }

    public override void OnLayerChange()
    {
        for (int i = 0; i < GetView().teamItems.Length; i++)
        {
            GameUtil.Instance.SetEffectLayer(GetView().teamItems[i]._obj_on.transform, _canvasOrder + 1);
        }
    }

    private int TeamIdx = -1;
    private void OnTeamItemClick(int idx)
    {
        TeamIdx = idx;
        OnBaseInfo();

        for (int i = 0; i < GetView().teamItems.Length; i++)
        {
            GetView().teamItems[i]._obj_on.SetActive(i == idx);
        }

        GetView().img_Pet.gameObject.SetActive(TeamIdx == 0);

        if (TeamIdx == 0)
            OnPetShow();

        ResouseHelper.Instance.LoadSpriteByAtlas("Common/Atlas/Role", TeamIdx == 0 ? "bg_juese" : "bg_huoban", GetView().img_Role);
    }

    //装备槽位信息
    private void OnBaseInfo()
    {
        string Str = "";
        if (TeamIdx == 0)
        {
            Str = observePlayerInfo.Nickname;
            if (observePlayerInfo.GuildName != null && !string.IsNullOrEmpty(observePlayerInfo.GuildName))
                Str = $"<size=24><color=#4DAE1F>[{observePlayerInfo.GuildName}]</color></size>\n{Str}";
        }
        else
        {
            PartnerSlotInfo solt = GetPartnerSlotInfo(TeamIdx);
            if (solt.PartnerUid != 0)
            {
                var cfg_Partner = CFG_Partner.GetById(GetPartnerInfo(solt.PartnerUid).CfgId);
                Str = GameUtil.Instance.GetFightStrByQuality(cfg_Partner.Quality, cfg_Partner.Name);
            }
        }

        GetView().txt_RoleName.text = Str;

        for (int i = 0; i < GetView().root_Equips.childCount; i++)
        {
            int idx = i + 1;
            var go = GetView().root_Equips.Find("equip_" + idx);

            SlotInfo equipInfo = TeamIdx == 0 ? observePlayerInfo.SlotInfo[idx] : GetPartnerSlotInfo(TeamIdx).Equips[idx];
            JUIExtensions.txt(go.Find("txt_Starlv")).text = equipInfo.Stars.ToString();
            go.Find("txt_Starlv").gameObject.SetActive(equipInfo.Stars > 0);
            JUIExtensions.txt(go.Find("txt_lv")).text = equipInfo.Enhance == 0 ? "" : "+" + equipInfo.Enhance;
           
            AddBoxItemBestSpine(equipInfo.Item, go.transform);

            for (int j = 0; j < equipInfo.Gemstone.Length; j++)
            {
                go.Find("inset").GetChild(j).gameObject.SetActive(equipInfo.Gemstone[j] != 0);
                if (equipInfo.Gemstone[j] != 0)
                {
                    GameUtil.Instance.SetItemIcon(equipInfo.Gemstone[j], JUIExtensions.img(go.Find("inset").GetChild(j)));
                }
            }
        }
        OnEquipShow();
    }

    //设置装备信息
    private void OnEquipShow()
    {
        string Str = "";
        string lvStr = "";
        if (TeamIdx == 0)
        {
            Str = GameUtil.Instance.FormatPowerWithUnit(observePlayerInfo.Combat);
            lvStr = observePlayerInfo.Level.ToString();

            string poolKey = "MainRole" + observePlayerInfo.Job + observePlayerInfo.Sex;
            SetRoleSpine(poolKey, GameConst.GetModelByJobAndSex((JobType)observePlayerInfo.Job, (SexType)observePlayerInfo.Sex));
        }
        else
        {
            var solt = GetPartnerSlotInfo(TeamIdx);
            string poolKey = "";
            ObservePlayerInfo.PartnerStats soltStats = GetPartnerState(solt.PartnerUid);
            if (solt.PartnerUid != 0 && soltStats != null)
            {
                Str = GameUtil.Instance.FormatPowerWithUnit(soltStats.Combat);
                PartnerInfo partnerInfo = GetPartnerInfo(solt.PartnerUid);
                lvStr = partnerInfo.Level.ToString();
                poolKey = "Partner_" + partnerInfo.CfgId;
            }

            SetRoleSpine(poolKey, poolKey == "" ? "" : GameConst.GetMonsterModelByID(poolKey));
        }

        GetView().txt_RolePower.text = Str;
        GetView().txt_RolePower.gameObject.SetActive(!string.IsNullOrEmpty(Str));

        GetView().txt_lv.text = lvStr;
        GetView().txt_lv.gameObject.SetActive(!string.IsNullOrEmpty(lvStr));


        for (int i = 0; i < GetView().root_Equips.childCount; i++)
        {
            int idx = i + 1;
            var go = GetView().root_Equips.Find("equip_" + idx);

            var itemInfo = TeamIdx == 0 ? observePlayerInfo.SlotInfo[idx].Item : GetPartnerSlotInfo(TeamIdx).Equips[idx].Item;
            JUIExtensions.img(go.Find("icon")).enabled = itemInfo != null;
            JUIExtensions.img(go.Find("iconbg")).enabled = itemInfo != null;
            JUIExtensions.btn(go).enabled = itemInfo != null;
            bool hasEquip = itemInfo != null;
            if (hasEquip)
                GameUtil.Instance.SetItemIcon(itemInfo.CfgId, JUIExtensions.img(go, "icon"), JUIExtensions.img(go, "iconbg"), null, (int)itemInfo.Equip.Quality);

            go.Find("defaultImg").gameObject.SetActive(!hasEquip);

            AddBoxItemBestSpine(itemInfo, go.transform);
        }
        SetArrtInfo(false);
    }

    private void AddBoxItemBestSpine(ItemInfo itemInfo, Transform item)
    {
        var bestEffectTrans = item.Find("BestEffect");
        var bestEffectGo = bestEffectTrans.gameObject;
        if (itemInfo != null && itemInfo.Equip != null && (int)itemInfo.Equip.Quality >= 4 && itemInfo.Equip.Attributes != null && itemInfo.Equip.Attributes.Count > 0)
        {
            bestEffectGo.SetActive(true);
            SpineUtil.Instance.PlayItemSpineLoop(bestEffectTrans, (int)itemInfo.Equip.Quality);
        }
        else
            bestEffectGo.SetActive(false);
    }

    private PartnerInfo GetPartnerInfo(long partnerUid)
    {
        return observePlayerInfo.PartnerInfo.Find(v => v.Uid == partnerUid);
    }

    private PartnerSlotInfo GetPartnerSlotInfo(int slotId)
    {
        return observePlayerInfo.Partners.Find(v => v.SlotId == slotId);
    }



    private bool AttrShow = false;

    //设置玩家属性信息
    private void SetArrtInfo(bool isShow)
    {
        AttrShow = isShow;
        if (TeamIdx == 0)
        {
            GetView().btn_Attribute.gameObject.SetActive(true);
            (GetView().txt_attr0.text, GetView().txt_attr1.text, GetView().txt_attr2.text, GetView().txt_attr3.text) = GetRoleAttribute((isShow ? -1 : 4), observePlayerInfo.Attributes);
        }
        else
        {
            var solt = GetPartnerSlotInfo(TeamIdx);
            GetView().btn_Attribute.gameObject.SetActive(solt.PartnerUid != 0);
            if (solt.PartnerUid != 0 && GetPartnerState(solt.PartnerUid) != null)
                (GetView().txt_attr0.text, GetView().txt_attr1.text, GetView().txt_attr2.text, GetView().txt_attr3.text) = GetRoleAttribute((isShow ? -1 : 4), GetPartnerState(solt.PartnerUid).Attributes);
        }

        GetView().arrow_Attribute.localRotation = Quaternion.Euler(0, 0, isShow ? 180 : 0);
        LayoutRebuilder.ForceRebuildLayoutImmediate(GetView().btn_Attribute.GetComponent<RectTransform>());
    }


    private void SetRoleSpine(string poolKey, string path)
    {

        if (GetView().root_Role.childCount > 0)
        {
            if (poolKey == GetView().root_Role.GetChild(0).gameObject.name)
                return;
        }
        void ClearRole()
        {
            int Count = GetView().root_Role.childCount;
            for (int i = 0; i < Count; i++)
            {
                PoolHelper.Instance.RecycObj(GetView().root_Role.GetChild(0).gameObject);
            }
        }
        if (!string.IsNullOrEmpty(poolKey))
        {
            _ = GameHelper.Instance.SetAsyncObj(poolKey, path, (_mainRole) =>
            {
                ClearRole();
                _mainRole.transform.SetParent(GetView().root_Role);
                _mainRole.transform.localScale = Vector3.one;
                _mainRole.transform.localPosition = Vector3.zero;
                _mainRole.SetActive(true);

                SpineUtil.Instance.PlayQueue(_mainRole.transform.Find("spine").GetComponent<SkeletonGraphic>(), new List<string> { "stand" }, null, true);
            });
        }
        else
            ClearRole();

    }

    private (string, string, string, string) GetRoleAttribute(int Count, List<Stats> StatsList)
    {
        int idx = 0;
        string Str0 = "";
        string Str1 = "";
        string Str2 = "";
        string Str3 = "";

        int Amount = 4;
        for (int i = 0; i < CFG_XMLConfig.xmlConfig.RolePanelList.Count; i++)
        {
            var item = CFG_XMLConfig.xmlConfig.RolePanelList[i];
            if (item.Show == 0)
                continue;

            if (idx % Amount == 0)
            {
                Str0 = $"{Str0}{item.Name}:<color=#ffffff>{CFG_XMLConfig.GetAttrValueToStr(item.Type, GetValue(item.SID))}</color>\n";
            }
            else if (idx % Amount == 1)
            {
                Str1 = $"{Str1}{item.Name}:<color=#ffffff>{CFG_XMLConfig.GetAttrValueToStr(item.Type, GetValue(item.SID))}</color>\n";
            }
            else if (idx % Amount == 2)
            {
                Str2 = $"{Str2}{item.Name}:<color=#ffffff>{CFG_XMLConfig.GetAttrValueToStr(item.Type, GetValue(item.SID))}</color>\n";
            }
            else
            {
                Str3 = $"{Str3}{item.Name}:<color=#ffffff>{CFG_XMLConfig.GetAttrValueToStr(item.Type, GetValue(item.SID))}</color>\n";
            }
            idx++;

            if (idx == Count)
                break;
        }

        long GetValue(int SId)
        {
            for (int i = 0; i < StatsList.Count; i++)
            {
                if (StatsList[i].attr == (Attribute)SId)
                    return StatsList[i].Num;
            }
            return 0;
        }

        return (Str0, Str1, Str2, Str3);
    }

    public override void OnExit()
    {

    }

    public override void OnDestroy()
    {

    }
}