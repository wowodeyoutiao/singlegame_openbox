using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Core;

public class UILoading : UIBase<UIModelLoading, UIViewLoading>
{
    public class Loading 
    {
        public float _time;
        public string _str;
        public int showType = 0; // 0:д╛хо
    }
    public override UILayer GetLayer()
    {
        return UILayer.Worning;
    }
    public override void OnInit()
    {

    }
    private TimeItem TimeJa;

    private void CloseTimeJa()
    {
        if (TimeJa != null)
        {
            GameHelper.Instance.RemoveInvoke(TimeJa);
            TimeJa = null;
        }
    }
    public override void OnEnter(params object[] args)
    {
        var data = args[0] as Loading;

        GetView().obj_Show0.SetActive(data.showType == 0);
        GetView().obj_Show1.SetActive(data.showType == 1);
        GetView().txt_tip.text = string.IsNullOrEmpty(data._str) ? "" : data._str;

        CloseTimeJa();

        if (data._time > 0)
        {
            TimeJa = GameHelper.Instance.AddInvoke(new TimeItem(data._time, () => {
                if (isNullGo)
                    return;
                UITool.Close<UILoading>(); 
            } ));
        }
    }

    public override void OnExit()
    {
        CloseTimeJa();
    }

    public override void OnDestroy()
    {
        CloseTimeJa();
    }
}