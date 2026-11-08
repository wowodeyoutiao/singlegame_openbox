using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using DG.Tweening;
using JEngine.Helper;
using JEngine.Core;

public class UIPowerTip : UIBase<UIModelPowerTip, UIViewPowerTip>
{
    public static bool IsOpne = false;
    public override UILayer GetLayer()
    {
        return UILayer.Tip;
    }

    public override void OnLayerChange()
    {
        GameUtil.Instance.SetEffectLayer(GetView().eff, _canvasOrder + 1);
    }

    public override void OnInit()
    {

    }

    private TimeItem TimeJa;
    private Tween tween;
    public override void OnEnter(params object[] args)
    {
        int nowNum = (int)args[0];
        int baseNum = (int)args[1];
        OnShow(nowNum, baseNum);
    }

    public void OnShow(int nowNum, int baseNum)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.UI_RolePowerUp);

        GetView()._Animator.SetTrigger("anim");
        ClearTimeJa();
        TimeJa = GameHelper.Instance.AddInvoke(new TimeItem(3, () => {

            UITool.Close<UIPowerTip>();
            ClearTimeJa();
        }));
        GetView().cGroup_add.alpha = 1;
        GetView().cGroup_add.DOKill();

        GetView().txt_base.text = "" + baseNum;
        GetView().txt_add.text = "+" + (nowNum - baseNum);
        GetView().cGroup_add.DOFade(0, 0.5f).SetDelay(1.5f);

        if (tween != null)
        {
            tween.Kill();
        }
        tween = DOTween.To(() => baseNum, value => GetView().txt_base.text = "" + value, nowNum, 2.5f);
    }

    private void ClearTimeJa()
    {
        if (GameHelper.Instance != null && TimeJa != null)
        {
            GameHelper.Instance.RemoveInvoke(TimeJa);
            TimeJa = null;
        }
    }


    public override void OnExit()
    {
        IsOpne = false;
        ClearTimeJa();
        tween = null;
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}