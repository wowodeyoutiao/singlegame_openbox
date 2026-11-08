using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using DG.Tweening;

public class UIJumpLoad : UIBase<UIModelJumpLoad, UIViewJumpLoad>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {

    }

    public override void OnEnter(params object[] args)
    {
        DOTween.To(() => 0, (a) => {
            GetView().txt_tip.text = a.ToString() + "%";
        }, 100, 3f).SetEase(Ease.Linear).OnComplete(() => {
            UITool.Close<UIJumpLoad>();
        });
    }

    public override void OnExit()
    {
        
    }

     public override void OnDestroy()
    {

    }
}