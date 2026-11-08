using HotUpdateScripts;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using MetaFramework.UI;

public class UIGameLuckCoin : UIBase<UIModelGameLuckCoin, UIViewGameLuckCoin>
{
    public override UILayer GetLayer()
    {
        return UILayer.Base;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIGameLuckCoin>();
        });

        GetView().btn_Coin.onClick.AddListener(() =>
        {
            if (_state != UIState.Start)
                return;
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            TcpC.C2S_FateCoin();
        });
    }

    private enum UIState
    {
        Start,
        Over
    }
    public void OnCoinOpne(int _cfgID)
    {
        cfgID = _cfgID;
        _state = UIState.Over;
        OnUIRefresh();
    }

    private UIState _state = UIState.Start;

    private void OnUIRefresh()
    {
        var CFGInfo = CFG_FateCoin.GetById(cfgID);
        string baseIcon = _state == UIState.Start || CFGInfo.EffectType == 3 ? "table_common_xgn_1" : CFGInfo.EffectType == 1 ? "table_common_xgn_2" : "table_common_xgn_3";
        ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CommonBGAtlas, baseIcon, GetView().img_Bg);

        baseIcon = _state == UIState.Start || CFGInfo.EffectType == 3 ? "bg_common_xgn_1" : CFGInfo.EffectType == 1 ? "bg_common_xgn_2" : "bg_common_xgn_3";
        ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CommonBGAtlas, baseIcon, GetView().img_zs1);
        ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CommonBGAtlas, baseIcon, GetView().img_zs2);

        baseIcon = _state == UIState.Start ? "word_myyb_zi" : CFGInfo.EffectType == 1 ? "word_xdy_zi" : CFGInfo.EffectType == 2 ? "word_ydb_zi" : "word_wsfs_zi";
        ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.FightAtlas, baseIcon, GetView().img_Title);

        GetView().txt_Tip.text = _state == UIState.Start ? "点击硬币开启命运的安排" : $"<color=#{(CFGInfo.EffectType == 1 ? "91e48d" : CFGInfo.EffectType == 2 ? "ff675e" : "fffbce")}>{CFGInfo.Des}</color>";

        GetView().root_Award.gameObject.SetActive(_state != UIState.Start);
        GetView().img_Coin.gameObject.SetActive(_state != UIState.Start);

        GetView().btn_Coin.gameObject.SetActive(_state == UIState.Start);
        if (_state != UIState.Start)
        {
            ClearTime();
            GetView().txt_Time.text = "点击空白区域关闭";
            ResouseHelper.Instance.LoadSpriteByAtlas(CFGInfo.Gallery, CFGInfo.Icon, GetView().img_Coin);
        }
        

        GetView().btn_Close.gameObject.SetActive(_state != UIState.Start);
    }

    private TimeItem TimeJaction;
    private void ClearTime()
    {
        if (TimeJaction != null && GameHelper.Instance)
        {
            GameHelper.Instance.RemoveInvoke(TimeJaction);
            TimeJaction = null;
        }
    }

    private int cfgID = 0;

    public override void OnEnter(params object[] args)
    {
        _state = UIState.Start;
        cfgID = 0;
        OnUIRefresh();
        int CDTime = CFG_XMLConfig.xmlConfig.BaseValue.FateCoinAutotime + 1;
        TimeJaction = GameHelper.Instance.AddInvoke(new TimeItem(CDTime, 1, () =>
        {
            CDTime--;
            GetView().txt_Time.text = CDTime + "秒后自动开启";
            if (CDTime <= 0)
            {
                TcpC.C2S_FateCoin();
                ClearTime();
            }
        }));
    }


    public override void OnExit()
    {
        ClearTime();
    }

    public override void OnDestroy()
    {
        ClearTime();
    }
}