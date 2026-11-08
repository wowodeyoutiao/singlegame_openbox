using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using MetaFramework.UI;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using DG.Tweening;

public class UIGuide : UIBase<UIModelGuide, UIViewGuide>
{
    private UIViewGuide guideView;
    private string norSpineName = "6666";
    private string dragSpineName = "oooo";
    private string comboSpineName = "9999";

    private CFGItem_Guidance cFGItem_Guidance;
    private Button targetBtn;
    private ExButton targetExBtn;
    private TouchEventCom targetTouchEventCom;

    private TimeItem guideTimeItem;

    private bool clickEnable = true;

    private Vector3 oldPos;

    public override UILayer GetLayer()
    {
        return UILayer.Guide;
    }

    public override void OnInit()
    {
        guideView = GetView();

        guideView.closeBtn.onClick.AddListener(OnClickCloseButton);
        guideView.skipBtn.onClick.AddListener(OnClickSkipButton);
        guideView.focusBtn.onClick.AddListener(OnClickFocusButton);
    }

    // 关闭
    private void OnClickCloseButton()
    {
        if(cFGItem_Guidance != null && cFGItem_Guidance.GuideType == (int)GuidanceGuideType.Text && cFGItem_Guidance.CompleteType == (int)GuidanceCompleteType.Clicked)
            FinishGuideAndDoNext();
        else 
            Debug.Log("检测当前引导是否结束");
    }

    // 跳过
    private void OnClickSkipButton()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

        var jumpToId = cFGItem_Guidance.JumpToId;
        GuidanceM.Instance.Jump(jumpToId);
        cFGItem_Guidance = GuidanceM.Instance.GetGuideConfig();
        if(cFGItem_Guidance == null)
        {
            GuidanceM.Instance.SaveGuide(jumpToId);
            UniGameEvent.Send(new GameEvent.Guide_Jump());
            UITool.Close<UIGuide>();
            return;
        }

        if (!ShowGuide())
            UITool.Close<UIGuide>();
    }

    // 点击小手
    private void OnClickFocusButton()
    {
        if (!clickEnable) return;

        if (cFGItem_Guidance.CompleteType != (int)GuidanceCompleteType.MainBattle)
            clickEnable = false;

        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

        //Debug.LogError($"{targetBtn}, {targetExBtn}, {targetTouchEventCom}");
        if (targetBtn != null)
            targetBtn.onClick.Invoke();

        if (targetExBtn != null)
            targetExBtn.OnClick.Invoke();

        if (targetTouchEventCom != null)
            targetTouchEventCom.OnPointerClick(null);

        // 完成条件（点击完成）
        if (cFGItem_Guidance.CompleteType == (int)GuidanceCompleteType.Clicked)
            FinishGuideAndDoNext();
    }

    public override void OnEnter(params object[] args)
    {
        UniGameEvent.AddEventListener<GameEvent.Refresh_Guide>(RefreshGuide);
        UniGameEvent.AddEventListener<GameEvent.Receive_Guide_ServerData>(ReceiveGuideServerData);

        if (!ShowGuide())
            UITool.Close<UIGuide>();
    }

    private bool ShowGuide()
    {
        cFGItem_Guidance = GuidanceM.Instance.GetGuideConfig();
        if (cFGItem_Guidance == null)
            return false;

        Debug.Log($"显示引导，ID: {cFGItem_Guidance.ID}, ClassParam:{cFGItem_Guidance.ClassParam}, TargetParam:{cFGItem_Guidance.TargetParam}");

        if(string.IsNullOrEmpty(cFGItem_Guidance.ClassParam))
        {
            ShowUI();
            return true;
        }

        // 打开UI
        if (!UITool.IsOpen(cFGItem_Guidance.ClassParam))
            DelayShowUI(0.5f);
        else
            DelayShowUI(0.2f);

        return true;
    }

    #region UI 
    private void ClearTimeItem()
    {
        if (guideTimeItem != null)
        {
            if (GameHelper.Instance != null)
                GameHelper.Instance.RemoveInvoke(guideTimeItem);
            guideTimeItem = null;
        }
    }

    private void DelayShowUI(float delay)
    {
        ClearTimeItem();
        var className = cFGItem_Guidance.ClassParam;
        guideTimeItem = GameHelper.Instance.AddInvoke(new TimeItem(10, delay, () =>
        {
            if (UITool.IsOpen(className))
            {
                ClearTimeItem();
                ShowUI();
            }
        }));
    }

    private void ShowUI()
    {
        if (isNullGo) return;

        // 完成条件（时间到）
        TimeOut();
        targetBtn = null;

        // 跳过
        guideView.skipGo.SetActive(cFGItem_Guidance.JumpToId != 0);

        guideView.maskGo.SetActive(cFGItem_Guidance.HalfMask == 0);
        guideView.battleMaskGo.SetActive(cFGItem_Guidance.HalfMask == 1);

        // 位置
        var pos = GetPosition();
        guideView.focusTrans.position = pos;
        pos = guideView.focusTrans.localPosition;
        pos.z = 0;
        guideView.focusTrans.localPosition = pos;

        guideView.tipGo.SetActive(false);

        // 小手是否出现
        var isShow = cFGItem_Guidance.GuideType >= (int)GuidanceGuideType.Click;
        if (isShow)
        {
            guideView.focusTrans.localPosition = oldPos;

            var enable = cFGItem_Guidance.GuideType != (int)GuidanceGuideType.BattleETCJoystick;
            if (cFGItem_Guidance.GuideType == (int)GuidanceGuideType.BattleETCJoystick)
                SpineUtil.Instance.PlayLoop(guideView.focusTrans, dragSpineName);
            else
            {
                if (cFGItem_Guidance.GuideParam == 1)
                    SpineUtil.Instance.PlayLoop(guideView.focusTrans, comboSpineName);
                else
                    SpineUtil.Instance.PlayLoop(guideView.focusTrans, norSpineName);
            }

            guideView.focusGo.SetActive(true);
            guideView.focusBtn.enabled = false;
            guideView.skeletonGraphic.raycastTarget = false;

            var dis = Vector3.Distance(oldPos, pos);
            if (dis < 100)
            {
                guideView.focusBtn.enabled = enable;
                guideView.focusTrans.localPosition = pos;
                guideView.skeletonGraphic.raycastTarget = enable;
                DelayShowText();
            }
            else
            {
                guideView.focusTrans.DOLocalMove(pos, 0.5f).OnComplete(() =>
                {
                    guideView.focusTrans.DOKill();
                    guideView.focusBtn.enabled = enable;
                    guideView.skeletonGraphic.raycastTarget = enable;
                    DelayShowText();
                });
            }

            oldPos = pos;
        }
        else
        {
            guideView.focusGo.SetActive(false);
            DelayShowText();
        }

        void DelayShowText()
        {
            clickEnable = true;

            // 文本
            isShow = cFGItem_Guidance.GuideType == (int)GuidanceGuideType.Text ||
                (cFGItem_Guidance.GuideType >= (int)GuidanceGuideType.TextAndClick && cFGItem_Guidance.GuideType != (int)GuidanceGuideType.BattleETCJoystick);
            guideView.tipGo.SetActive(isShow);
            guideView.desText.text = cFGItem_Guidance.Describe;
            if (isShow)
            {
                if (pos.x + 350 >= Screen.width * 0.5f)
                {
                    pos.x -= 180;
                    guideView.leftArrowGo.SetActive(true);
                    guideView.rightArrowGo.SetActive(false);
                }
                else
                {
                    pos.x += 240;
                    guideView.leftArrowGo.SetActive(false);
                    guideView.rightArrowGo.SetActive(true);
                }
                pos.z = 0;
                guideView.tipGo.transform.localPosition = pos;
            }
        }
    }

    private Vector3 GetPosition()
    {
        if (string.IsNullOrEmpty(cFGItem_Guidance.TargetParam))
            return Vector3.zero;

        Transform targetTrans = null;
        switch (cFGItem_Guidance.ClassParam)
        {
            case "UIHall":
                {
                    var view = UITool.Get<UIHall>().GetView();
                    targetTrans = view.transform.Find(cFGItem_Guidance.TargetParam);
                    
                }
                break;

            case "UIRole":
                {
                    var view = UITool.Get<UIRole>().GetView();
                    targetTrans = view.transform.Find(cFGItem_Guidance.TargetParam);
                }
                break;

            case "UIPartnerCall":
                {
                    var view = UITool.Get<UIPartnerCall>().GetView();
                    targetTrans = view.transform.Find(cFGItem_Guidance.TargetParam);
                }
                break;

            case "UIAwards":
                {
                    var view = UITool.Get<UIAwards>().GetView();
                    targetTrans = view.transform.Find(cFGItem_Guidance.TargetParam);
                }
                break;

            case "UIPartnerUp":
                {
                    var view = UITool.Get<UIPartnerUp>().GetView();
                    targetTrans = view.transform.Find(cFGItem_Guidance.TargetParam);
                }
                break;

            case "UIExpPool":
                {
                    var view = UITool.Get<UIExpPool>().GetView();
                    targetTrans = view.transform.Find(cFGItem_Guidance.TargetParam);
                }
                break;

            case "UIFightChoice":
                {
                    var view = UITool.Get<UIFightChoice>().GetView();
                    targetTrans = view.transform.Find(cFGItem_Guidance.TargetParam);
                }
                break;

            case "UIGameFightOver":
                {
                    var view = UITool.Get<UIGameFightOver>().GetView();
                    targetTrans = view.transform.Find(cFGItem_Guidance.TargetParam);
                }
                break;

            case "UIVIP":
                {
                    var view = UITool.Get<UIVIP>().GetView();
                    targetTrans = view.transform.Find(cFGItem_Guidance.TargetParam);
                }
                break;

            case "UIGameDungeons":
                {
                    var view = UITool.Get<UIGameDungeons>().GetView();
                    targetTrans = view.transform.Find(cFGItem_Guidance.TargetParam);
                }
                break;

            case "UIGameOverAwards":
                {
                    var view = UITool.Get<UIGameOverAwards>().GetView();
                    targetTrans = view.transform.Find(cFGItem_Guidance.TargetParam);
                }
                break;
        }

        if(targetTrans != null)
        {
            if (targetTrans != null)
            {
                targetBtn = targetTrans.GetComponent<Button>();
                targetExBtn = targetTrans.GetComponent<ExButton>();
                targetTouchEventCom = targetTrans.GetComponent<TouchEventCom>();

                var rectTransform = targetTrans.GetComponent<RectTransform>();
                Vector3[] corners = new Vector3[4];
                rectTransform.GetWorldCorners(corners);

                // 中心点 = 左下角和右上角的中点
                return new Vector3((corners[0].x + corners[2].x) / 2, (corners[0].y + corners[2].y) / 2, (corners[0].z + corners[2].z) / 2);
            }
            else
                Debug.LogError($"查找不到ID:{cFGItem_Guidance.ID},ClassParam: {cFGItem_Guidance.ClassParam}，TargetParam:{cFGItem_Guidance.TargetParam}, 看看配置对不对");
        }

        return Vector3.zero;
    }
    #endregion

    #region 完成条件（时间到）
    private void TimeOut()
    {
        if (cFGItem_Guidance.CompleteType != (int)GuidanceCompleteType.TimeOut)
            return;

        if(cFGItem_Guidance.CompleteValue == 0)
            FinishGuideAndDoNext();
        else
        {
            GameHelper.Instance.SetDeayTimeDo(() =>
            {
                FinishGuideAndDoNext();
            }, cFGItem_Guidance.CompleteValue * 0.001f);
        }
    }
    #endregion

    #region 完成条件（任务、主线）
    private void RefreshGuide(GameEvent.Refresh_Guide args)
    {
        if(cFGItem_Guidance == null)
        {
            Debug.Log("cFGItem_Guidance is null");
            return;
        }

        if (cFGItem_Guidance.CompleteType != (int)args.completeType)
            return;

        Debug.Log($"判断引导是否完成: {args.completeType}, {args.value}");
        if(args.completeType == GuidanceCompleteType.Task)
        {
            if (TaskM.Instance.CheckTaskFinish(cFGItem_Guidance.CompleteValue))
                FinishGuideAndDoNext();
        }
        else if(args.completeType == GuidanceCompleteType.MainBattle)
        {
            if(args.value >= cFGItem_Guidance.CompleteValue)
                FinishGuideAndDoNext();
        }
        else if (args.completeType == GuidanceCompleteType.Battle)
        {
            if (args.value >= cFGItem_Guidance.CompleteValue)
                FinishGuideAndDoNext();
        }
        else if (args.completeType == GuidanceCompleteType.ReceiveServer)
        {
            if (args.value >= cFGItem_Guidance.CompleteValue)
                FinishGuideAndDoNext();
        }
    }
    #endregion

    #region 结束当前引导，开始下一关
    private void FinishGuideAndDoNext()
    {
        if (isNullGo) return;

        guideView.skipGo.SetActive(false);
        guideView.focusGo.SetActive(false);
        guideView.tipGo.SetActive(false);
        targetBtn = null;
        targetExBtn = null;
        targetTouchEventCom = null;
        guideView.maskGo.SetActive(true);

        // 上报数数
        SendToShuShu();

        if (cFGItem_Guidance.NoticeServer == 1)
            GuidanceM.Instance.SaveGuide(cFGItem_Guidance.ID);
        else
        {
            if (cFGItem_Guidance.DelayTime > 0 && GameHelper.Instance != null)
            {
                GameHelper.Instance.SetDeayTimeDo(() =>
                {
                    DoNext();
                }, cFGItem_Guidance.DelayTime * 0.001f);
            }
            else
                DoNext();
        }
    }

    private void SendToShuShu()
    {
        Dictionary<string, object> properties = new Dictionary<string, object>();
        properties.Add("name", "guide_completed");
        properties.Add("guide_step", $"{cFGItem_Guidance.ID}");
        properties.Add("guide_name", cFGItem_Guidance.Describe);
        TDAnalyticsUtil.Track("guide_completed", properties, true);
    }

    private void ReceiveGuideServerData(GameEvent.Receive_Guide_ServerData args)
    {
        Debug.Log($"保存引导ID，收到服务器返回数据：{args.Code}, {args.Error}, {args.Guidance[0]}");
        if(args.Code != 0)
        {
            UITool.ShowTips(args.Error);
            return;
        }

        if (cFGItem_Guidance == null)
            return;

        var delayTime = cFGItem_Guidance.DelayTime;
        cFGItem_Guidance = null;
        if (delayTime > 0 && GameHelper.Instance != null)
        {
            GameHelper.Instance.SetDeayTimeDo(() =>
            {
                DoNext();
            }, delayTime * 0.001f);
        }
        else
            DoNext();
    }

    private void DoNext()
    {
        Debug.Log("开始引导下一步");
        GuidanceM.Instance.MoveNext();

        if (!ShowGuide())
            UITool.Close<UIGuide>();
    }
    #endregion

    public override void OnExit()
    {
        ClearTimeItem();
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_Guide>(RefreshGuide);
        UniGameEvent.RemoveEventListener<GameEvent.Receive_Guide_ServerData>(ReceiveGuideServerData);
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}