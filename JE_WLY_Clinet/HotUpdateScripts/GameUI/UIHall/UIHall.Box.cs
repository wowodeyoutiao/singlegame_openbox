using Cysharp.Threading.Tasks;
using DG.Tweening;
using HotUpdateScripts;
using HotUpdateScripts.ProtoM;
using JEngine.Core;
using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;
using MetaFramework.UI;
using pb;
using UnityEngine;

public partial class UIHall
{
    private bool isRequiredOpenBox = false;
    private bool isAutoClothing = false; // 正在穿戴
    private bool isShowBoxUI; // 是否显示弹窗

    // 计时器
    private TimeItem AutoOpenBoxTimeItem;
    private TimeItem waitOpenBoxResult;
    private TimeItem waitPupopCompareEquip;


    private bool isClickBoxHandleBtn;

    private void OnInitHallBox()
    {
        // 一键回收
        uIViewHall.boxDelBtn.onClick.AddListener(() =>
        {
            OnClickBoxDelButton();
        });

        // 一键穿戴
        uIViewHall.boxClothingBtn.onClick.AddListener(() =>
        {
            OnClickBoxClothingButton();
        });

        uIViewHall.boxSpineBtn.onClick.AddListener(OnClickBoxButton);
    }

    private void OnClickBoxDelButton()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        long[] uids = new long[RoleM.Instance.OpenBoxItems.Count];
        for (int i = 0; i < RoleM.Instance.OpenBoxItems.Count; i++)
        {
            uids[i] = RoleM.Instance.OpenBoxItems[i].Uid;
        }
        isClickBoxHandleBtn = true;
        RecycleItems(uids);
    }

    private void OnClickBoxClothingButton()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        isClickBoxHandleBtn = true;
        OnClothingEquip();
    }

    #region 开箱子
    private void CheckBox()
    {
        isRequiredOpenBox = false;
        uIViewHall.boxNodeTrans.gameObject.SetActive(hallBattleState == HallBattleState.Rest && RoleM.Instance.TreasuerBox > 0);
        uIViewHall.boxNumText.text = RoleM.Instance.TreasuerBox.ToString();
    }

    private void OnClickBoxButton()
    {
        if (isRequiredOpenBox) return;

        if (hallBattleState != HallBattleState.Rest|| RoleM.Instance.TreasuerBox == 0)
            return;

        SpineUtil.Instance.PlayOnce(uIViewHall.boxSpineBtn.transform, "dead", () =>
        {
            CheckBox();
            SpineUtil.Instance.PlayLoop(uIViewHall.boxSpineBtn.transform, "stand");
        });

        if (RoleM.Instance.OpenBoxItems.Count > 0)
        {
            Debug.Log("当前有展示的宝箱");
            return;
        }

        if (!BagM.Instance.CheckEquipVacancy(CFG_XMLConfig.xmlConfig.BaseValue.Offlinebox))
        {
            UITool.ShowTips("背包空间置不足，无法开箱。请清理背包装备后再试。");
            return;
        }

        if (RoleM.Instance.TreasuerBox <= 0)
        {
            UITool.ShowTips("开启次数不足");
            return;
        }

        isRequiredOpenBox = true;
        int openLevel = RoleM.Instance.GetDefaultOpenBoxEquipLevel();
        int openJob = RoleM.Instance.GetDefaultOpenBoxEquipJob();

        TcpC.C2S_HomeOpenBox(1, openLevel, openJob);
    }
    #endregion


    private void RefreshHallBox()
    {
        if (RoleM.Instance.OpenBoxItems != null)
            OnRefreshHomeBox(default);
    }

    // 穿戴装备
    private void OnClothingEquip()
    {
        isAutoClothing = true;
        uIViewHall.boxClothingBtn.gameObject.SetActive(false);
        ClothingEquip();
    }

    private void ClothingEquip()
    {
        ItemInfo item = GetOpenFirstUpEquip();
        if (item == null)
            return;
        (_, int soltIdx) = GameUtil.Instance.GetEquipIsBest(item);
        TcpC.C2S_ChangeEquip(item.Uid, 1, (int)CFG_Item.GetById(item.CfgId).Equiptype, soltIdx);
    }

    private void ClearAutoClothDelayTime()
    {
        if (GameHelper.Instance != null && autoClothDelay != null)
        {
            GameHelper.Instance.RemoveInvoke(autoClothDelay);
        }
    }

    [Subscriber(ThreadMode.Main)]
    private void OnSync_HomeOpenBox(GameEvent.Sync_HomeOpenBox _e)
    {
        if (isNullGo) return;

        if (_e.data.Code == 0)
        {
            waitOpenBoxResult = GameHelper.Instance.AddInvoke(new TimeItem(0.8f, () =>
            {
                if (isNullGo) return;
                OnRefreshHomeBox(default);
            }));
        }
    }

    //宝箱开启刷新
    private void OnRefreshHomeBox(GameEvent.Refresh_HomeOpen _e)
    {
        uIViewHall.oneBoxOpenGo.gameObject.SetActive(RoleM.Instance.OpenBoxItems.Count == 1);
        uIViewHall.multipleBoxOpenTrans.gameObject.SetActive(RoleM.Instance.OpenBoxItems.Count > 1);

        if (RoleM.Instance.OpenBoxItems.Count > 1)
        {
            for (int i = 0; i < uIViewHall.multipleBoxOpenTrans.childCount; i++)
            {
                var root = uIViewHall.multipleBoxOpenTrans.GetChild(i);
                root.gameObject.SetActive(i < RoleM.Instance.OpenBoxItems.Count);
                if (i < RoleM.Instance.OpenBoxItems.Count)
                {
                    var item = RoleM.Instance.OpenBoxItems[i];
                    GameUtil.Instance.SetItemIcon(item.CfgId, JUIExtensions.img(root, "icon"), JUIExtensions.img(root), null, (int)item.Equip.Quality);
                    GameUtil.Instance.SetItemJobIcon(item.CfgId, JUIExtensions.img(root, "icon_job"));

                    AddBoxItemBestSpine(item, root);

                    (int powerNum, int soltIdx) = GameUtil.Instance.GetEquipIsBest(item);
                    root.Find("up").gameObject.SetActive(powerNum > 0 && soltIdx >= 0);
                    root.Find("down").gameObject.SetActive(powerNum < 0 && soltIdx >= 0);
                    JUIExtensions.btn(root).onClick.RemoveAllListeners();

                    JUIExtensions.btn(root).onClick.AddListener(() =>
                    {
                        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                        // 正在一键穿装备
                        if (isAutoClothing)
                            return;
                        
                        UITool.Open<UIItemUseOrDel>(null, item, true);
                        if (soltIdx >= 0 && soltIdx != curRoleIndex)
                            OnEquipTabBtnClick(soltIdx);
                    });
                }
            }
        }
        else if (RoleM.Instance.OpenBoxItems.Count == 1)
        {
            var item = RoleM.Instance.OpenBoxItems[0];
            GameUtil.Instance.SetItemJobIcon(item.CfgId, JUIExtensions.img(uIViewHall.oneBoxOpenGo.transform, "icon_job"));
            GameUtil.Instance.SetItemIcon(item.CfgId, JUIExtensions.img(uIViewHall.oneBoxOpenGo.transform, "icon"), JUIExtensions.img(uIViewHall.oneBoxOpenGo.transform), null, (int)item.Equip.Quality);

            AddBoxItemBestSpine(item, uIViewHall.oneBoxOpenGo.transform);

            (int powerNum, int soltIdx) = GameUtil.Instance.GetEquipIsBest(item);
            uIViewHall.oneBoxOpenGo.transform.Find("up").gameObject.SetActive(powerNum > 0 && soltIdx >= 0);
            uIViewHall.oneBoxOpenGo.transform.Find("down").gameObject.SetActive(powerNum < 0 && soltIdx >= 0);

            JUIExtensions.btn(uIViewHall.oneBoxOpenGo.transform).onClick.RemoveAllListeners();
            JUIExtensions.btn(uIViewHall.oneBoxOpenGo.transform).onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                // 正在一键穿装备
                if (isAutoClothing)
                    return;

                if (!GuidanceM.Instance.isGuiding)
                    UITool.Open<UIItemUseOrDel>(null, item, true);
                if (soltIdx >= 0 && soltIdx != curRoleIndex)
                    OnEquipTabBtnClick(soltIdx);
            });
        }

        uIViewHall.boxClothingBtn.gameObject.SetActive(GetOpenUpEquipCount() > 0 && !isAutoClothing);

        if (GetOpenUpEquipCount() > 0 && isAutoClothing)
        {
            autoClothDelay = GameHelper.Instance.AddInvoke(new TimeItem(0.07f, () =>
            {
                if (!isNullGo) ClothingEquip();
                ClearAutoClothDelayTime();
            }));
        }

        isShowBoxUI = RoleM.Instance.OpenBoxItems.Count > 0;

        if (GetOpenUpEquipCount() == 0 && isAutoClothing)
        {
            if(isShowBoxUI && autoHallBattle)
                OnClickBoxDelButton();

            isAutoClothing = false;
            UniGameEvent.Send(new GameEvent.Refresh_NewPlayerGuide() { args = new object[1] { "RecycleFinish" } });
        }

        uIViewHall.boxShowUI.gameObject.SetActive(isShowBoxUI);
        uIViewHall.boxShowUIMaskGo.transform.SetAsLastSibling();
        uIViewHall.boxShowUIMaskGo.SetActive(isShowBoxUI && hallBattleState == HallBattleState.Rest);

        if (RoleM.Instance.OpenBoxItems.Count == 0 && hallBattleState != HallBattleState.Rest && isClickBoxHandleBtn)
        {
            isClickBoxHandleBtn = false;
            MoveToNext();
        }
    }


    #region 启用自动开箱
   
    private void AddBoxItemBestSpine(ItemInfo itemInfo, Transform item)
    {
        var bestEffectTrans = item.Find("BestEffect");
        var bestEffectGo = bestEffectTrans.gameObject;
        if (itemInfo.Equip != null && (int)itemInfo.Equip.Quality >= 4 && itemInfo.Equip.Attributes != null && itemInfo.Equip.Attributes.Count > 0)
        {
            bestEffectGo.SetActive(true);
            SpineUtil.Instance.PlayItemSpineLoop(bestEffectTrans, (int)itemInfo.Equip.Quality);
        }
        else
            bestEffectGo.SetActive(false);
    }

    // 回收所有开箱开出的道具或装备
    private void RecycleItems(long[] uids)
    {
        TcpC.C2S_ItemDisassembly(uids);
    }

    private void DisassemblyEffectFlyAsync(GameEvent.OpenBoxDisassemblyEquip _e)
    {
        if (_e.count == 0) return;

        AudioManager.Instance.PlaySound(AudioKey.Sound.UI_FlyEffect);
        for (int i = 0; i < _e.count; i++)
        {
            LoadDisassemblyEffect().Forget();
        }
    }

    private async UniTask LoadDisassemblyEffect()
    {
        if (!PoolHelper.Instance.CheckPool("UI_ExpUp"))
        {
            var baseObj = new JPrefab(GameConst.GetEffPath("UI_ExpUp"), true);
            await UniTask.WaitUntil(() => baseObj.Loaded);
            PoolHelper.Instance.AddPool("UI_ExpUp", baseObj.Instance);
        }
        await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool("UI_ExpUp"));
        GameObject effObj = PoolHelper.Instance.GetPoolObj("UI_ExpUp");
        effObj.transform.SetParent(GetView().effectRoot);
        effObj.transform.localScale = Vector3.one;
        effObj.transform.localPosition = GetRandomPos();
        effObj.SetActive(true);
        float duration = UnityEngine.Random.Range(0.6f, 1.2f);
        effObj.transform.DOLocalMove(downUI.GetButtonPositionRelativeCanvas(6), duration).SetEase(Ease.InOutQuad).OnComplete(() =>
        {
            PoolHelper.Instance.RecycObj(effObj);
        });

        Vector3 GetRandomPos()
        {
            var pos = uIViewHall.boxShowUI.transform.localPosition;
            var rect = uIViewHall.boxShowUI.transform as RectTransform;
            pos.y = pos.y + rect.rect.height / 2 + 35;
            float yrange = (rect.rect.height - 110) / 2;
            return new Vector3(UnityEngine.Random.Range(pos.x - 140f, pos.x + 140f), UnityEngine.Random.Range(pos.y - yrange, pos.y + yrange));
        }
    }
    #endregion

    #region 开箱
    private int GetOpenUpEquipCount()
    {
        int count = 0;
        for (int i = 0; i < RoleM.Instance.OpenBoxItems.Count; i++)
        {
            var item = RoleM.Instance.OpenBoxItems[i];
            (int powerNum, int soltIdx) = GameUtil.Instance.GetEquipIsBest(item);
            if (powerNum > 0 && soltIdx >= 0)
            {
                count++;
            }
        }
        return count;
    }

    private ItemInfo GetOpenFirstUpEquip()
    {
        for (int i = 0; i < RoleM.Instance.OpenBoxItems.Count; i++)
        {
            var item = RoleM.Instance.OpenBoxItems[i];
            (int powerNum, int soltIdx) = GameUtil.Instance.GetEquipIsBest(item);
            if (powerNum > 0 && soltIdx >= 0)
            {
                return item;
            }
        }
        return null;
    }
    #endregion

    #region 删除计时器
    private void ClearAutoOpenTime()
    {
        if (GameHelper.Instance != null && AutoOpenBoxTimeItem != null)
        {
            GameHelper.Instance.RemoveInvoke(AutoOpenBoxTimeItem);
            AutoOpenBoxTimeItem = null;
        }
    }

    private void ClearTimeOpenResult()
    {
        if (GameHelper.Instance != null && waitOpenBoxResult != null)
        {
            GameHelper.Instance.RemoveInvoke(waitOpenBoxResult);
            waitOpenBoxResult = null;
        }
    }

    private void ClearTimeCompareEquip()
    {
        if (GameHelper.Instance != null && waitPupopCompareEquip != null)
        {
            GameHelper.Instance.RemoveInvoke(waitPupopCompareEquip);
            waitPupopCompareEquip = null;
        }
    }
    #endregion
}