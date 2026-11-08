using MetaFramework.UI;
using UnityEngine;
using UnityEngine.UI;

public class UIViewStrengthen : MonoBehaviour, IUIView
{
    public Image grayImg;
    [HideInInspector]
    public Material grayMat;
    public Button closeBtn;

    public Transform btnTrans;

    public void Awake()
    {
        grayMat = grayImg.material;
    }
    public void Init(GameObject go)
    {

    }

    public void OnDestroy()
    {
        UITool.RemoveUI<UIStrengthen>();
    }
}