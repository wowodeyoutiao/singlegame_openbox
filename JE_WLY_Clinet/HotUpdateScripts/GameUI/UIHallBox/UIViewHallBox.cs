using MetaFramework.UI;
using UnityEngine;
using UnityEngine.UI;

public class UIViewHallBox : MonoBehaviour, IUIView
{
    public Button lvBtn;
    public Text lvText;
    public GameObject lvRedGo;
    public Text numText;
    public Text openNumText;
    public Button leftBtn;
    public Button rightBtn;
    public Button settingBtn;
    public Button atuoBtn;
    public Transform autoTrans;
    public Button okBtn;

    public Transform boxTrans;

    public void Awake()
    {
    }

    public void Init(GameObject go)
    {
    }

    public void OnDestroy()
    {
        UITool.RemoveUI<UIHallBox>();
    }
}