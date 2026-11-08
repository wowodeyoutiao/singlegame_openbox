using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using DG.Tweening;
using JEngine.Helper;
using JEngine.UI;
using HotUpdateScripts;
using System.Collections.Generic;
using Spine.Unity;
using HotUpdateScripts.ProtoM;

public class UIGameOverAwards : UIBase<UIModelGameOverAwards, UIViewGameOverAwards>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        var closeBtn = GetView().obj_CloseTip.GetComponent<Button>();
        closeBtn.onClick.AddListener(() =>
        {
            if (_showState == ShowState.ShowOneAward)
                OnShowAllAward();
            else if (_showState == ShowState.ShowAllAward)
                UITool.Close<UIGameOverAwards>();
        });

        GetView().btn_Mask.onClick.AddListener(() =>
        {
            if (_showState == ShowState.ShowOneAward)
                OnShowAllAward();
            else if (_showState == ShowState.ShowAllAward)
                UITool.Close<UIGameOverAwards>();
        });

        for (int i = 0; i < GetView().root_Cards.childCount; i++)
        {
            int idx = i;
            GetView().root_Cards.GetChild(i).GetComponent<Button>().onClick.AddListener(() => 
            {
                if (_showState != ShowState.StarAnim) return;
                AudioManager.Instance.PlaySound(AudioKey.Sound.UI_SelectCard);
                OnResultBtnClick(idx);
            });
            GetView().root_Cards.GetChild(i).Find("Item").GetComponent<Button>().onClick.AddListener(() => {
                AudioManager.Instance.PlaySound(AudioKey.Sound.UI_SelectCard);
                UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = CFGAwardInfo[idx].id });
            });
        }
    }

    private int _SelectIdx = 0;
    private void OnResultBtnClick(int idx = 0)
    {
        if (_showState != ShowState.StarAnim) return;

        UniGameEvent.Send(new GameEvent.Refresh_NewPlayerGuide() { args = new object[1] { "BeginShowCardPlayEffect" } });

        _SelectIdx = idx;
        _showState = ShowState.ShowOneAward;
        if (idx != selectIdx)
        {
            var temp = CFGAwardInfo[idx];
            CFGAwardInfo[idx] = CFGAwardInfo[selectIdx];
            CFGAwardInfo[selectIdx] = temp;
        }

        int bestItemIndex = -1;

        for (int i = 0; i < GetView().root_Cards.childCount; i++)
        {
            var cardRoot = GetView().root_Cards.GetChild(i);
            //var itemRoot = cardRoot.Find("Item");
            if (_SelectIdx == i)
            {
                int _idx = i;
                cardRoot.DOLocalMoveY(20, 0.1f).SetLoops(2, LoopType.Yoyo).OnComplete(() => {
                    //itemRoot.gameObject.SetActive(true);
                    cardRoot.Find("select").gameObject.SetActive(true);
                    var spineRoot = cardRoot.Find("spine");
                    spineRoot.gameObject.SetActive(true);
                    if (CFGAwardInfo[_idx].rate == MinRate)
                    {
                        SpineUtil.Instance.PlayQueue(spineRoot.GetComponent<SkeletonGraphic>(), new List<string> { "rong", "sao" }, () => {
                        }, true);
                        bestItemIndex = i;
                    }
                    else
                        SpineUtil.Instance.PlayOnce(spineRoot, "rong", () =>
                        {
                            spineRoot.gameObject.SetActive(false);
                        });

                    ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CardAtlas, CFG_Item.GetById(CFGAwardInfo[_idx].id).CardIcon, JUIExtensions.img(cardRoot.Find("bg")));

                    JUIExtensions.txt(cardRoot.Find("txt_num")).text = $"{GameUtil.Instance.GetItemName(CFGAwardInfo[_idx].id)}*{GameUtil.Instance.FormatNumWithUnit(CFGAwardInfo[_idx].num)}" ;
                    GameHelper.Instance.SetDeayTimeDo(OnShowAllAward, 1f);
                });
               
            }

        }

        if (bestItemIndex == idx)
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.UI_ShowCard);
        }
        OnTimeClear();

        if (GuidanceM.Instance.isGuiding)
            UniGameEvent.Send(new GameEvent.Refresh_Guide { completeType = GuidanceCompleteType.ReceiveServer, value = 0 });
    }

    private void OnShowAllAward()
    {
        if (_showState != ShowState.ShowOneAward) return;

        UniGameEvent.Send(new GameEvent.Refresh_NewPlayerGuide() { args = new object[1] { "EndShowCardPlayEffect" } });

        _showState = ShowState.ShowAllAward;
        for (int i = 0; i < GetView().root_Cards.childCount; i++)
        {
            var cardRoot = GetView().root_Cards.GetChild(i);
            int _idx = i;
            if (i != _SelectIdx)
            {
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CardAtlas, CFG_Item.GetById(CFGAwardInfo[_idx].id).CardIcon, JUIExtensions.img(cardRoot.Find("bg")));
                JUIExtensions.txt(cardRoot.Find("txt_num")).text = $"{GameUtil.Instance.GetItemName(CFGAwardInfo[_idx].id)}*{GameUtil.Instance.FormatNumWithUnit(CFGAwardInfo[_idx].num)}";

                var spineRoot = cardRoot.Find("spine");
                spineRoot.gameObject.SetActive(true);
                if (CFGAwardInfo[_idx].rate == MinRate)
                    SpineUtil.Instance.PlayQueue(spineRoot.GetComponent<SkeletonGraphic>(), new List<string> { "rong", "sao" }, () =>
                    {
                    }, true);
                else
                    SpineUtil.Instance.PlayOnce(spineRoot, "rong", () => {
                        spineRoot.gameObject.SetActive(false);
                    });
            }
        }
        GetView().obj_CloseTip.SetActive(true);
    }

    private CFGItem_Battle.AwardItem[] CFGAwardInfo;
    private int selectIdx = 0;

    private int MinRate = 0;

    private enum ShowState {No,MoveAnim ,StarAnim,ShowOneAward, ShowAllAward, ShowResult }
    private ShowState _showState = ShowState.StarAnim;
    public override void OnEnter(params object[] args)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.UI_RandomCard);

        var _CFGAwardInfo = CFG_Battle.GetById((int)args[0]).Randomaward;
        CFGAwardInfo = new CFGItem_Battle.AwardItem[_CFGAwardInfo.Length];
        _CFGAwardInfo.CopyTo(CFGAwardInfo, 0);
        MinRate = _CFGAwardInfo[0].rate;
        for (int i = 0; i < _CFGAwardInfo.Length; i++)
        {
            if (_CFGAwardInfo[i].id == (int)args[1] && _CFGAwardInfo[i].num == (int)args[2])
                selectIdx = i;

            if (_CFGAwardInfo[i].rate < MinRate && _CFGAwardInfo[i].rate > 0)
                MinRate = _CFGAwardInfo[i].rate;
        }

        GetView().obj_CloseTip.SetActive(false);
        _showState = ShowState.No;
        OnAnimStar();

        NewPlayerGuideWidget.Instance.AddGuideGameObject("GameOverAward_Card", GetView().root_Cards);
        NewPlayerGuideWidget.Instance.AddGuideGameObject("GameOverAward_CloseBtnPos", GetView().obj_CloseTip.transform);
        NewPlayerGuideWidget.Instance.AddGuideGameObject("GameOverAward_CloseBtn", GetView().btn_Mask.transform);
    }

    private void OnAnimStar()
    {
        for (int i = 0; i < GetView().root_Cards.childCount; i++)
        {
            var cardRoot = GetView().root_Cards.GetChild(i);
            cardRoot.gameObject.SetActive(false);
            cardRoot.localPosition = Vector3.zero;
            cardRoot.Find("spine").gameObject.SetActive(false);
            cardRoot.Find("select").gameObject.SetActive(false);
            JUIExtensions.txt(cardRoot.Find("txt_num")).text = string.Empty;

            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.FightAtlas, "btn_cards_bm", JUIExtensions.img(cardRoot.Find("bg")));
        }
        var BaseRoot = GetView().root_Cards.GetChild(0);
        BaseRoot.gameObject.SetActive(true);
        BaseRoot.localPosition = Vector3.down * 1000;

        DOTween.Sequence().Append(BaseRoot.DOLocalMoveY(0, 0.5f))
            .Insert(0f, BaseRoot.DOScaleX(0.1f, 0.2f))
            .Insert(0.2f, BaseRoot.DOScaleX(1, 0.2f)).OnComplete(() => {

                for (int i = 0; i < GetView().root_Cards.childCount; i++)
                {
                    var cardRoot = GetView().root_Cards.GetChild(i);
                    cardRoot.gameObject.SetActive(true);
                    cardRoot.DOLocalMoveX(-210 + i * 210, 0.2f).SetEase(Ease.Linear).SetDelay(0.1f).OnComplete(() => {
                        if (_showState == ShowState.No)
                            DoMultiSwap();
                    });
                }
               
            });

      
    }

    private List<int[]> moveIds = new() { new int[] { 2, 0, 1 }, new int[] { 1, 2, 0 }, new int[] { 0, 1, 2 }, new int[] { 2, 0, 1 }, new int[] { 1, 2, 0 }, new int[] { 0, 1, 2 } };

    private void DoMultiSwap()
    {
        _showState = ShowState.MoveAnim;
        for (int i = 0; i < moveIds.Count; i++)
        {
            int idx = i;
            GameHelper.Instance.SetDeayTimeDo(() =>
            {
                SwapPair(moveIds[idx], idx == moveIds.Count - 1);
            }, idx * 0.25f);
        }
    }

    private void SwapPair(int[] idxs,bool isLast = false)
    {
        for (int i = 0; i < GetView().root_Cards.childCount; i++)
        {
            GetView().root_Cards.GetChild(i).DOLocalMoveX(-210 + idxs[i] * 210, 0.2f).SetEase(Ease.Linear).OnComplete(() => {
                if (_showState == ShowState.MoveAnim && isLast)
                {
                    _showState = ShowState.StarAnim;
                    OnTimeStar();
                    UniGameEvent.Send(new GameEvent.Refresh_NewPlayerGuide() { args = new object[1] { "OpenGameOverAwardUI" } });
                }
            });
        }
    }

    private void OnTimeStar()
    {
        int Count = 11;
        OnTimeClear();
        if(!GuidanceM.Instance.isGuiding)
        {
            Ja = GameHelper.Instance.AddInvoke(new TimeItem(Count, 1, () =>
            {
                if (isNullGo) return;
                Count--;
                GetView().txt_Time.text = Count + "秒后自动翻牌";
                if (Count <= 0)
                {
                    OnResultBtnClick();
                }
            }));
        }else
            GetView().txt_Time.text = string.Empty;
    }

    TimeItem Ja;
    private void OnTimeClear()
    {
        if (Ja != null && GameHelper.Instance != null)
        {
            GameHelper.Instance.RemoveInvoke(Ja);
            Ja = null;
        }
        GetView().txt_Time.text = string.Empty;
    }


    public override void OnExit()
    {
        
    }

    public override void OnDestroy()
    {

    }
}