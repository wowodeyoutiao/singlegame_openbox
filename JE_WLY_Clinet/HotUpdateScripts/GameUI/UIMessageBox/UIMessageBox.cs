using MetaFramework.UI;
using System;

public class UIMessageBox : UIBase<UIModelMessageBox, UIViewMessageBox>
{
    public override UILayer GetLayer()
    {
        return UILayer.Top;
    }

    public class CallBack 
    {
        public bool isTogg;
        public string ShowStr;
        public Action YesFun;
        public Action NoFun;
        public string BtnYes;
        public string BtnNo;

        public Action TogYesFun;
        public Action TogNoFun;
    }

    private CallBack callBack;

    public override void OnInit()
    {
        GetView().Btn_Yes.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIMessageBox>();
            callBack.YesFun?.Invoke();
        });
        GetView().Btn_No.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIMessageBox>();
            callBack.NoFun?.Invoke();
        });

        GetView().toggle_tip.onValueChanged.AddListener((isOn) => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (isOn)
                callBack.TogYesFun?.Invoke();
            else
                callBack.TogNoFun?.Invoke();
        });
    }

    public override void OnEnter(params object[] args)
    {
        callBack = args[0] as CallBack;

        GetView().Text_des.text = callBack.ShowStr;

        GetView().Text_Yes.text = !string.IsNullOrEmpty(callBack.BtnYes) ? callBack.BtnYes : "确定";
        GetView().Text_No.text = !string.IsNullOrEmpty(callBack.BtnNo) ? callBack.BtnNo : "取消";

        GetView().toggle_tip.gameObject.SetActive(callBack.isTogg);
        if (callBack.isTogg)
            GetView().toggle_tip.isOn = false;
    }

    public override void OnExit()
    {
    }
}