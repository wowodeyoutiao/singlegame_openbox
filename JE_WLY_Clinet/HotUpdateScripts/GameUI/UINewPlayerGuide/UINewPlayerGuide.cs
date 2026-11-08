using DG.Tweening;
using HotUpdateScripts.ProtoM;
using ILRuntime.Runtime.Debugger.Protocol;
using JEngine.Helper;
using MetaFramework.UI;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class UINewPlayerGuide : UIBase<UIModelNewPlayerGuide, UIViewNewPlayerGuide>
{
    private NewPlayerGuideStepBase  curStep = null;
    public override UILayer GetLayer()
    {
        return UILayer.Guide;
    }

    public override void OnInit()
    {
        
    }

    public override void OnEnter(params object[] args)
    {
        if (GetView().guideSteps == null || args.Length < 1)
            OnDestroy();


        var guideStep = (eGuideStep)args[0];

        if (!GetView().guideSteps.ContainsKey(guideStep))
            OnDestroy();

        if (curStep != null)
        {
            curStep?.gameObject.SetActive(false);
        }

        var test = GetView().guideSteps;

        var stepContoller = GetView().guideSteps[guideStep];

        curStep = stepContoller;
        curStep?.gameObject.SetActive(true);
        curStep?.Enter(this);
    }

    public override void OnExit()
    {
        //ClearTimeItem();
        //UniGameEvent.RemoveEventListener<GameEvent.Refresh_Guide>(RefreshGuide);
        //UniGameEvent.RemoveEventListener<GameEvent.Receive_Guide_ServerData>(ReceiveGuideServerData);
    }

    public override void OnDestroy()
    {
        curStep?.Exit();

        OnExit();
    }
}