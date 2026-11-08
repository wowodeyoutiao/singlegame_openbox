using DG.Tweening;
using HotUpdateScripts;
using HotUpdateScripts.ProtoM;
using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;
using MetaFramework.UI;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using static CFG_XMLConfig.XML_Config;

public class UIPet : UIBase<UIModelPet, UIViewPet>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public Camera mainCam;
    public override void OnInit()
    {
        mainCam = GameObject.Find("UICamera").GetComponent<Camera>();
        GetView().btn_Close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIPet>();
        });

        GetView().btn_help.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIComRule>(null,43,"抽取概率公示");
        });

        GetView().btn_help2.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIComRule>(null,42);
        });

        GetView().btn_PetBook.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIPetBook>();
        });

        GetView().btn_res.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = (int)CFG_Item.ItemID.PetChangeCard });
        });

        GetView().btn_PetGift.onClick.AddListener(() =>
        {
            if (GameM.Instance.IsAudit)
                UITool.ShowTips("功能暂未开放！");
            else
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                UITool.Open<UIGiftBag>(null, ActivityType.PetGift);
            }
        });

        GetView().btn_Yc.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIPetUp>();
        });

        GetView().touch_Pet.AddClickListener((e) =>
        {
            if (GetView().obj_Petlv.activeSelf)
            {
                var petInfo = PetM.Instance.GetOnePetByIdx(0);
                if (petInfo.CfgId == 0)
                    return;
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                UITool.Open<UIPetShow>(null, petInfo.CfgId, petInfo.Level, petInfo.Star);
            }
        });

        GetView().touch_Pet.AddDragBeginListener((p) =>
        {
            OnDragBegin(p, 0, GetView().img_Pet.gameObject);
        });
        GetView().touch_Pet.AddDragListener((p) => { 
            OnDrag(p, 0); 
        });
        GetView().touch_Pet.AddDragEndListener((p) => { 
            OnDragEnd(p, 0,true); 
        });
        GetView().touch_Pet.AddPointerUpListener((p) => { 
            OnDragEnd(p, 0, false);
        });

        for (int i = 0; i < GetView().PetItems.Length; i++)
        {
            int idx = i;
            GetView().PetItems[i].touchCom.AddDragBeginListener((p) =>
            {
                OnDragBegin(p, (idx + 1), GetView().PetItems[idx].img_Icon.gameObject);
            });
            GetView().PetItems[i].touchCom.AddDragListener((p) => { OnDrag(p, (idx + 1)); });
            GetView().PetItems[i].touchCom.AddDragEndListener((p) => { OnDragEnd(p, (idx + 1),true); });
            GetView().PetItems[i].touchCom.AddPointerUpListener((p) => { OnDragEnd(p, (idx + 1),false); });

            GetView().PetItems[i].touchCom.AddClickListener((e) =>
            {

                var petInfo = PetM.Instance.GetOnePetByIdx((idx + 1));
                if (petInfo.CfgId == 0)
                    return;
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                UITool.Open<UIPetShow>(null, petInfo.CfgId, petInfo.Level, petInfo.Star);
            });
        }

        FuncOpenShowWidget.Instance.HideFuncNewFlag(new List<SwitchOpne>() { SwitchOpne.Pet });
    }

    public override void OnLayerChange()
    {
        if (GetView().obj_fly != null)
            GameUtil.Instance.SetEffectLayer(GetView().obj_fly.transform, _canvasOrder + 1);
        else
        {
            _ = GameHelper.Instance.SetAsyncObj("UIFly", GameConst.GetEffPath("UIFly"), (go) => {
                go.transform.SetParent(uiGo.transform);
                go.gameObject.SetActive(false);
                go.transform.localScale = Vector3.one * 0.8f;
                GetView().obj_fly = go;
                GameUtil.Instance.SetEffectLayer(GetView().obj_fly.transform, _canvasOrder + 1);
            }, true);
        }

        for (int i = 0; i < GetView().root_Btns.childCount; i++)
        {
            GameUtil.Instance.SetEffectLayer(GetView().root_Btns.GetChild(i), _canvasOrder + 1);
        }
    }

    //所有解锁槽位的宠物数量
    private int AllPetCount = 0;
    //宠物槽位解锁
    private void OnPetLock()
    {
        AllPetCount = 0;
        for (int i = 0; i < GetView().PetItems.Length; i++) {

            int SoltIdx = i+1;
            bool isLock = CFG_XMLConfig.xmlConfig.PetBagDic.ContainsKey(SoltIdx) && CFG_XMLConfig.xmlConfig.PetBagDic[SoltIdx] > RoleM.Instance.Level;
            GetView().PetItems[i].obj_lock.SetActive(isLock);
            if (isLock)
                GetView().PetItems[i].txt_lock.text = "角色" + CFG_XMLConfig.xmlConfig.PetBagDic[SoltIdx] + "级解锁";
            else
                AllPetCount ++ ;
        }

        JUIExtensions.txt(GetView().root_Btns.GetChild(1), "consume/txt_usedNum").text = GameUtil.Instance.FormatNumWithUnit(CFG_XMLConfig.xmlConfig.BaseValue.Petbuy_needmoney);
        GameUtil.Instance.SetItemIcon(CFG_XMLConfig.xmlConfig.BaseValue.Petbuy_itemid, JUIExtensions.img(GetView().root_Btns.GetChild(1), "consume/img"));
    }

    private int BaseSolt = -1;
    private GameObject BaseImg = null;
    private void OnDragBegin(UnityEngine.EventSystems.PointerEventData pEvent, int idx, GameObject baseObj)
    {
        BaseSolt = idx;
        var petInfo = PetM.Instance.GetOnePetByIdx(idx);
        if (petInfo.CfgId == 0)
            return;
        AllCanDoPetShow(true);
        GetView().img_Drag.gameObject.SetActive(true);
        BaseImg = baseObj;
        BaseImg.SetActive(false);
        GetView().img_Drag.transform.position = baseObj.transform.position;

        GameUtil.Instance.SetItemIcon(CFG_Pet.GetById(petInfo.CfgId).ItemID, GetView().img_Drag);

    }

    private void AllCanDoPetShow(bool isShow)
    {
        if (BaseSolt < 0 || PetM.Instance.GetOnePetByIdx(BaseSolt).CfgId == 0)
            return;
        var pets = PetM.Instance.GetPets();
        var basePet = CFG_Pet.GetById(PetM.Instance.GetOnePetByIdx(BaseSolt).CfgId);
        for (int i = 0; i < AllPetCount; i++)
        {
            int soltIdx = i + 1;
            if (pets[soltIdx].CfgId != 0 && BaseSolt != soltIdx)
            {
                if (isShow)
                {
                    if (basePet.Mark == 0)
                        GetView().PetItems[i].obj_NoMove.SetActive(false);
                    else if (basePet.Mark == 1 && basePet.Quality == 8)
                        GetView().PetItems[i].obj_NoMove.SetActive(true);
                    else
                        GetView().PetItems[i].obj_NoMove.SetActive(CFG_Pet.GetById(pets[soltIdx].CfgId).Quality != basePet.Quality || CFG_Pet.GetById(pets[soltIdx].CfgId).Mark == 0);
                }
                else
                {
                    GetView().PetItems[i].obj_NoMove.SetActive(false);
                }
            }
        }
    }

    private void OnDrag(UnityEngine.EventSystems.PointerEventData pEvent, int idx)
    {
        var petInfo = PetM.Instance.GetOnePetByIdx(idx);
        if (petInfo.CfgId == 0)
            return;

        var dragPos = mainCam.ScreenToWorldPoint(pEvent.position);
        GetView().img_Drag.transform.position = new Vector3(dragPos.x, dragPos.y, GetView().img_fly.transform.position.z);
    }

    private void OnDragEnd(UnityEngine.EventSystems.PointerEventData pEvent, int idx,bool isEnd)
    {
        var petInfo = PetM.Instance.GetOnePetByIdx(idx);
        if (petInfo.CfgId == 0)
        {
            return;
        }

        GetView().img_Drag.gameObject.SetActive(false);

        if (BaseImg != null)
        {
            BaseImg.SetActive(true);
            BaseImg = null;
        }

        AllCanDoPetShow(false);
        GameObject target = pEvent.pointerCurrentRaycast.gameObject;

        if (target != null && target.tag == "skill" && isEnd)
        {
            int tagsolt = int.Parse(target.gameObject.name.Replace("PetItem_", ""));

            //原來位置 和锁定位置不处理
            if (tagsolt != BaseSolt && tagsolt <= AllPetCount )
            {
                if (PetM.Instance.GetOnePetByIdx(tagsolt).CfgId != 0)
                {
                    var basePet = CFG_Pet.GetById(PetM.Instance.GetOnePetByIdx(BaseSolt).CfgId);
                    var targetPet = CFG_Pet.GetById(PetM.Instance.GetOnePetByIdx(tagsolt).CfgId);
                    if (tagsolt != 0)
                    {
                        if (basePet.Quality != targetPet.Quality && basePet.Mark != 0 && targetPet.Mark != 0)
                        {
                            UITool.ShowTips("不同品质无法合成");
                        }
                    }
                }
                    
                TcpC.C2S_PetMove(BaseSolt, tagsolt);

                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "角色",
                    component = "宠物",
                    position = "",
                    event_click = "合成宠物",
                });
            }
        }
    }
    private void OnBindRed()
    {
        RedPointManager.Instance.RPS.Register(RedPointConst.Pet_UpStar.ToString(), (node) =>
        {
            if (isNullGo) return;
            GetView().obj_redDot.gameObject.SetActive(node.pointCount > 0);
        });
    }

    private CFGItem_Summon[] CFGItems;
    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<Event_Source_Refresh>(Event_Refresh);
        UniGameEvent.AddEventListener<GameEvent.Sync_TakeSource>(Sync_TakeSource);

        OnBindRed();
        OnPetLock();
        GetHttp();
        CFGItems = CFG_Summon.GetPetItems();
        OnEggBtnNumShow(true);
        OnPetsShow();
        OnTimeShow();

        for (int i = 0; i < GetView().root_Btns.childCount; i++)
        {
            int soltIdx = i;
            var btn = GetView().root_Btns.GetChild(i).GetComponent<ExButton>();
            btn.ClearAllListeners();
            btn.OnClick = () =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.UI_PetEggCall);
                OnEggBtnClick(soltIdx);
            };

            btn.OnPress = () =>
            {
                OnEggBtnClick(soltIdx,true);
            };
        }
    }

    private HTTP_PrizedrawListResponse PetrawList;

    //宠物蛋展示
    private void OnEggBtnNumShow(bool isInit = false)
    {
        for (int i = 0; i < GetView().root_Btns.childCount; i++) 
        {
            long itemCount = BagM.Instance.GetNumByCfgID(CFGItems[i].Consume[0].id);
            JUIExtensions.txt(GetView().root_Btns.GetChild(i), "txt_num").text = GameUtil.Instance.FormatNumWithUnit(itemCount);
            GetView().root_Btns.GetChild(i).Find("dot").gameObject.SetActive(BagM.Instance.GetNumByCfgID(CFGItems[i].Consume[0].id) >0);
            //if (CFGItems[i].Drawtype == 3 && itemCount == 0 && isShowUI )
            //{
            //    isShowUI = false;
            //    UITool.Open<UIActive>(null, ActivityType.PetGift);
            //}
        }

        if (isInit)
            GameUtil.Instance.SetItemIcon((int)CFG_Item.ItemID.PetChangeCard, GetView().img_res);
        GetView().txt_res.text = GameUtil.Instance.FormatNumWithUnit(BagM.Instance.GetNumByCfgID((int)CFG_Item.ItemID.PetChangeCard));
    }

    private bool isClick = false;
    private void OnEggBtnClick(int soltIdx,bool isPress = false)
    {
        if (isClick)
        {
            UITool.ShowTips("操作频繁！");
            return;
        }

        int canNum = 0;
        for (int i = 0; i < PetrawList.canplay.Length; i++)
        {
            if (PetrawList.canplay[i].type == (soltIdx + 1))
            {
                canNum = PetrawList.canplay[i].can_play;
                break;
            }
        }
        if (canNum == 0)
        {
            UITool.ShowTips("每日上限抽取次数已满！");
            return;
        }
        var BoxNum = BagM.Instance.GetNumByCfgID(CFGItems[soltIdx].Consume[0].id);
        if (BoxNum == 0)
        {
            if (soltIdx == 0)
            {
                UITool.ShowTips("数量不足");
                return;
            }
            else if (soltIdx == 2)
            {
                if (isPress)
                    return;

                if(GameM.Instance.IsAudit)
                    UITool.ShowTips("数量不足，可通过活动获取。");
                else
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Close<UIPet>();
                    UITool.Open<UIGiftBag>(null, ActivityType.PetGift);
                }
                return;
            }
            else {

                string timeStr = GameUtil.Instance.GetDataTime(GameUtil.Instance.GetCurrentTimeSeconds()).DayOfYear.ToString();
                if (GameUtil.Instance.GetPlayerPrefs(GameConst.PlayerPrefsKey.PetBuy.ToString()) == timeStr)
                {
                    Click();

                }
                else
                    UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
                    {
                        ShowStr = $"是否消耗{CFG_XMLConfig.xmlConfig.BaseValue.Petbuy_needmoney}{GameUtil.Instance.GetItemName(CFG_XMLConfig.xmlConfig.BaseValue.Petbuy_itemid)}购买一只随机宠物？",
                        isTogg = true,
                        TogNoFun = () =>
                        {
                            GameUtil.Instance.SetPlayerPrefs(GameConst.PlayerPrefsKey.PetBuy.ToString(), "");
                        },
                        TogYesFun = () =>
                        {
                            GameUtil.Instance.SetPlayerPrefs(GameConst.PlayerPrefsKey.PetBuy.ToString(), timeStr);
                        },

                        YesFun = Click,
                    });
                return;
            }
        }
        void Click()
        {
            isClick = true;
            //GameHelper.Instance.SetDeayTimeDo(() => { isClick = false; }, 0.3f);

            int flyCount = 0;
            var pets = PetM.Instance.GetPets();
            for (int i = 1; i < pets.Count; i++)
            {
                if (pets[i].CfgId == 0)
                {
                    flyCount = i;
                    break;
                }
            }

            HttpC.GetPetReceive(1, (soltIdx + 1), (_r) => {
                isClick = false;
                if (_r.code == 0)
                {
                    GetHttp();
                    if (flyCount == 0) return;
                    GetView().obj_fly.transform.position = GetView().root_Btns.GetChild(soltIdx).position;
                    GetView().obj_fly.SetActive(true);
                    GetView().obj_fly.transform.DOMove(GetView().PetItems[(flyCount - 1)].txt_lock.transform.position, 0.3f).OnComplete(() =>
                    {
                        GetView().obj_fly.SetActive(false);
                    });
                }

                string eggcolor = soltIdx switch
                {
                    0 => "绿色",
                    1 => "紫色",
                    2 => "橙色",
                    _ => ""
                };
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "角色",
                    component = "宠物",
                    position = "",
                    event_click = $"召唤{eggcolor}宠物蛋",
                });
            }, () => {
                isClick = false;
            });
        }
        Click();

    }

    //宠物展示
    private void OnPetsShow()
    {
        var pets = PetM.Instance.GetPets();
        //上阵宠物位
        bool isHavePet = pets[0].CfgId != 0;
        if (isHavePet)
        {
            //GameUtil.Instance.SetItemIcon(CFG_Pet.GetById(pets[0].CfgId).ItemID, GetView().img_Pet);
            OnPetRole(pets[0].CfgId, GetView().img_Pet.transform,0.7f);
            GetView().txt_Petlv.text =  "lv" + pets[0].Level;
            GetView().img_PetMark.enabled = CFG_Pet.GetById(pets[0].CfgId).Mark != 0;
            GetView().obj_Exp.SetActive(CFG_Pet.GetById(pets[0].CfgId).Mark == 0);
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CommonAtlas, "pet_" + CFG_Pet.GetById(pets[0].CfgId).Mark, GetView().img_PetMark);
        }

        GetView().obj_Exp.SetActive(isHavePet && CFG_Pet.GetById(pets[0].CfgId).Mark == 0);
        GetView().txt_PetName.text = isHavePet ? GameUtil.Instance.GetStrByQuality(CFG_Pet.GetById(pets[0].CfgId).Quality, CFG_Pet.GetById(pets[0].CfgId).Name) : "";
        GetView().img_Pet.gameObject.SetActive(isHavePet);
        GetView().obj_Petlv.SetActive(isHavePet);
        GameUtil.Instance.SetItemBg(isHavePet ? CFG_Pet.GetById(pets[0].CfgId).Quality : 0, JUIExtensions.img(GetView().touch_Pet.transform));

        //宠物列表
        for (int i = 0; i < AllPetCount; i++)
        {
            isHavePet = pets[i+1].CfgId != 0;
            GetView().PetItems[i].img_Icon.gameObject.SetActive(isHavePet);
            GetView().PetItems[i].obj_lv.SetActive(isHavePet);
            if (isHavePet)
            {
                var cdgPet = CFG_Pet.GetById(pets[(i + 1)].CfgId);
                OnPetRole(pets[i + 1].CfgId, GetView().PetItems[i].img_Icon.transform);
                GetView().PetItems[i].txt_lv.text = "lv" + pets[(i + 1)].Level;
                GetView().PetItems[i].img_Pet.enabled = cdgPet.Mark != 0;
                GetView().PetItems[i].obj_Exp.SetActive(cdgPet.Mark == 0);
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CommonAtlas, "pet_" + cdgPet.Mark, GetView().PetItems[i].img_Pet);
            }
            else
                GetView().PetItems[i].obj_Exp.SetActive(false);

            GameUtil.Instance.SetItemBg(isHavePet ? CFG_Pet.GetById(pets[(i + 1)].CfgId).Quality : 0, JUIExtensions.img(GetView().PetItems[i].touchCom.transform));
        }
    }

    private void OnPetRole(int cfgId,Transform root,float scal = 0.6f)
    {
        string pooName = "Pet_" + cfgId;
        _ = GameHelper.Instance.SetAsyncObj(pooName, GameConst.GetMonsterModelByID(pooName), (_mainRole) =>
        {
            int Count = root.childCount;
            for (int i = 0; i < Count; i++)
            {
                PoolHelper.Instance.RecycObj(root.GetChild(0).gameObject);
            }
            _mainRole.transform.SetParent(root);
            _mainRole.transform.localScale = Vector3.one * scal;
            _mainRole.transform.localPosition = new Vector3(0,- 40,0);
            _mainRole.SetActive(true);
        });
    }

    private TimeItem EggTime;
    private void ClearEggTime()
    {
        if (EggTime != null)
        { 
            GameHelper.Instance.RemoveInvoke(EggTime);
            EggTime = null;
        }
    }

    private void OnTimeShow()
    {
        ClearEggTime();
        int needTime = (int)(PetM.Instance.PetRecoveryTime - GameUtil.Instance.GetCurrentTimeSeconds());
        bool isShowTime = needTime > 0 && BagM.Instance.GetNumByCfgID(CFGItems[0].Consume[0].id) < CFG_XMLConfig.xmlConfig.BaseValue.Petfree_maxNum;
        GetView().obj_Time.SetActive(isShowTime);
        if (isShowTime)
        {
            EggTime = GameHelper.Instance.AddInvoke(new TimeItem(needTime, 1, () => { 
                GetView().txt_Time.text = GameUtil.Instance.GetTimeSToHMS(needTime);
                needTime--;
                if (needTime <= 0)
                {
                    GetView().obj_Time.SetActive(false);
                    ClearEggTime();
                }
            }));
        }
    }


    private void GetHttp()
    {
        HttpC.GetPrivilegeDetail(2, (_d) =>
        {
            PetrawList = _d;
            for (int i = 0; i < GetView().root_Btns.childCount; i++)
            {
                JUIExtensions.txt(GetView().root_Btns.GetChild(i), "txt_Limit").text = "每日上限:" + PetM.Instance.CanGetNums[i];
            }
            isClick = false;
        });
    }

    [Subscriber(ThreadMode.Main)]
    private void Refresh_PetMove(GameEvent.Refresh_PetMove _e)
    {
        if (_e.idx == -1)
            return;
        _ = GameHelper.Instance.SetAsyncObj("UIBoom", GameConst.GetEffPath("UIBoom"), (go) =>
        {
            go.transform.SetParent(_e.idx == 0 ? GetView().img_Pet.transform : GetView().PetItems[_e.idx - 1].img_Icon.transform);
            go.gameObject.SetActive(true);
            GameUtil.Instance.SetEffectLayer(go.transform, _canvasOrder + 1);
            go.transform.localScale = Vector3.one * 30;
            go.transform.localPosition = Vector3.zero;

            GameHelper.Instance.SetDeayTimeDo(() =>
            {
                if (isNullGo) return;
                PoolHelper.Instance.RecycObj(go);
            }, 0.5f);
        }, true);

        if (_e.petCfgId == 0)
            return;
        OnUIPetShow(_e.petCfgId);
    }

    private void OnUIPetShow(int petCfgId)
    {
        var CFGPet = CFG_Pet.GetById(petCfgId);
        if (CFGPet.Mark == 1 && CFGPet.Quality >= 7)
            UITool.Open<UIPetAward>(null, petCfgId);
    }

    [Subscriber(ThreadMode.Main)]
    private void Refresh_PetsInfo(GameEvent.Refresh_PetsInfo _e)
    {
        OnPetsShow();
        OnEggBtnNumShow();
        OnTimeShow();
    }

    //[Subscriber(ThreadMode.Main)]
    private void Event_Refresh(Event_Source_Refresh _e)
    {
        OnEggBtnNumShow();
        OnTimeShow();
    }

    //[Subscriber(ThreadMode.Main)]
    private void Sync_TakeSource(GameEvent.Sync_TakeSource _e)
    {
        OnEggBtnNumShow();
    }

    public override void OnExit()
    {
        CFGItems = null;
        ClearEggTime();
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<Event_Source_Refresh>(Event_Refresh);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_TakeSource>(Sync_TakeSource);
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}