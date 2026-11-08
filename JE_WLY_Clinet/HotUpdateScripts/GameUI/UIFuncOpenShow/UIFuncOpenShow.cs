using BestHTTP.SecureProtocol.Org.BouncyCastle.Asn1.X509;
using BestHTTP.SecureProtocol.Org.BouncyCastle.Ocsp;
using DG.Tweening;
using HotUpdateScripts;
using HotUpdateScripts.ProtoM;
using JEngine.Event;
using JEngine.Helper;

using MetaFramework.UI;
using System;
using System.Collections.Generic;
using UnityEngine;
using static CFG_XMLConfig.XML_Config;

public class UIFuncOpenShow : UIBase<UIModelFuncOpenShow, UIViewFuncOpenShow>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    
    public override void OnInit()
    {
        //GetView().btn_Click.onClick.AddListener(() =>
        //{
        //    if (seq != null) return;//等待动画播放完成继续下一个
        //    GetView().btn_Click.enabled = false;
        //    GetView().CancelInvoke();
        //    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        //    //seq?.Kill();
        //    onPlayEffect();

        //    GetView().btn_Click.gameObject.SetActive(!(curFuncList.Count == 0));
        //});
    }

    private List<int> curFuncList = null;//当前功能展示列表
    private XML_OpenFunction curFuncData = null; //当前展示功能数据
    private SwitchOpne curTag = SwitchOpne.No;
    private CanvasGroup canvasGroup = null;
    //private Tween rotateEffect;
    private Sequence seq;

    public override void OnEnter(params object[] args)
    {
        //GetView().btn_Click.gameObject.SetActive(true);
        //rotateEffect = GetView().root_Effect.DOLocalRotate(new Vector3(0, 0, 90), 5f, RotateMode.FastBeyond360).SetEase(Ease.Linear);
        //rotateEffect.SetLoops(-1, LoopType.Incremental);

        canvasGroup = GetView().tf_Icon.GetComponent<CanvasGroup>();

        curFuncList = (List<int>)args[0];
        //curFuncList.Add(curFuncList[0]);
        tryPlayEffect();
    }

    private void tryPlayEffect()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.UI_FunctionUnlock);

        //GetView().btn_Click.enabled = true;
        if (curFuncList == null || curFuncList.Count == 0)
        {
            UITool.Close<UIFuncOpenShow>();
            //触发新手引导。。。
            //curFuncData.
            return;
        }

        //Debug.LogError("尝试播放功能开启效果！");
        curFuncData = null;
        curTag = (SwitchOpne)curFuncList[0];

        if (CFG_XMLConfig.xmlConfig.SwitchOpenFunction.ContainsKey(curTag))
            curFuncData = CFG_XMLConfig.xmlConfig.SwitchOpenFunction[curTag];
        //Debug.LogError("新功能名称："+ curFuncData.name+" Icon:"+curFuncData.Icon+"Enum:"+ curTag); 
        if (curFuncData == null)
        {
            UITool.Close<UIFuncOpenShow>();
            return;
        }

        TcpC.C2S_RecordOpenID(curFuncList[0]);   //移除记录功能

        curFuncList.RemoveAt(0);
        GetView().SetAutoAction(onPlayEffect);

        //还原原始信息
        GetView().tf_Icon.localPosition = Vector3.zero;
        GetView().tf_Icon.localScale = Vector3.one;
        canvasGroup.alpha = 1f;
        GetView().goTitleRoot.SetActive(true);

        ResouseHelper.Instance.LoadSpriteByAtlas("Common/Atlas/FuncOpenShow", curFuncData.Icon, GetView().img_Icon);
    }

    private void onPlayEffect()
    {
        //Debug.LogError("开始表演");
        var target = FuncOpenShowWidget.Instance.GetFlyToTarget(curTag);

        if (target == null)
        {
            //Debug.LogError("未找到目标对象！，尝试播放下一个");
            tryPlayEffect();
            return;
        }

        seq = DOTween.Sequence();
        Tween tweenMove = GetView().tf_Icon.DOMove(target.position, 1f);
        tweenMove.SetEase(Ease.OutCubic);

        seq.AppendCallback(() => GetView().goTitleRoot.SetActive(false));

        seq.Join(GetView().tf_Icon.DOScale(Vector3.zero, 1.5f));
        seq.Join(tweenMove);
        seq.Join(canvasGroup.DOFade(0f, 1));

        seq.AppendInterval(0.1f); // 等0.1秒
        seq.AppendCallback(() =>
        {
            seq.Kill();
            seq = null;
            //Debug.LogError("表演结束，尝试播放下一个");
            tryPlayEffect();
        });
    }

    public override void OnDestroy()
    {
        canvasGroup = null;
        curFuncList = null;//当前功能展示列表
        curFuncData = null; //当前展示功能数据
        curTag = SwitchOpne.No;
        canvasGroup = null;

        //rotateEffect?.Kill();
    }

    public override void OnLayerChange()
    {
        int order = this._canvasOrder;
        ParticleSystemRenderer[] particleSystems = GetView().root_Effect.GetComponentsInChildren<ParticleSystemRenderer>();
        for (int i = 0; i < particleSystems.Length; i++)
        {
            particleSystems[i].sortingOrder = order + i;
        }
    }
}