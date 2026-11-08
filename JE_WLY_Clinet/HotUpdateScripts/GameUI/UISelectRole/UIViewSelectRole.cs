using MetaFramework.UI;
using UnityEngine;
using UnityEngine.UI;

public class UIViewSelectRole : MonoBehaviour, IUIView
{
    public Transform Root_Btn;
    public Image img_job;
    public Text txt_JobName;
    public Button Btn_Creat;
    public Button Btn_Delete;
    public Button Btn_StarGame;
    public Transform root_Role;

    public Transform root_JobDes;
    public Text[] JobDesTxts;

    public Button btn_Back;

    public void Awake()
    {
        JobDesTxts = new Text[4];
        for (int i = 0; i < 4; i++)
        {
            JobDesTxts[i] = root_JobDes.Find("txt_Des" + i).GetComponent<Text>();
        }
    }
    public void Init(GameObject go)
    {

    }
    private void OnDestroy()
    {
        UITool.RemoveUI<UISelectRole>();
    }
}