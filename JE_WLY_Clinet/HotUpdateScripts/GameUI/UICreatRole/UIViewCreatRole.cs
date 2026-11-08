using MetaFramework.UI;
using UnityEngine;
using UnityEngine.UI;

public class UIViewCreatRole : MonoBehaviour, IUIView
{
    public Transform root_role;
    public Transform root_JobSelect;
    public Button[] JobSelectBtns;
    public Transform root_Select;

    public Image img_Boy;
    public Image img_Gril;
    public Image img_Role;

    public InputField input_Name;
    public Button btn_randomName;
    public Button btn_Creat;

    public Text txt_Job;
    public Image img_job;

    public Transform root_job;
    public Text[] JobDesTxts;
    public Button btn_back;

    public void Awake()
    {
        JobSelectBtns = new Button[6];
        for (int i = 0; i < JobSelectBtns.Length; i++)
        {
            JobSelectBtns[i] = root_JobSelect.Find("Btn_" + (i + 1)).GetComponent<Button>();
        }

        JobDesTxts = new Text[4];
        for (int i = 0; i < 4; i++)
        {
            JobDesTxts[i] = root_job.Find("txt_Des" + i).GetComponent<Text>();
        }

    }
    public void Init(GameObject go)
    {

    }
    void OnDestroy()
    {
        UITool.RemoveUI<UICreatRole>();
    }
}