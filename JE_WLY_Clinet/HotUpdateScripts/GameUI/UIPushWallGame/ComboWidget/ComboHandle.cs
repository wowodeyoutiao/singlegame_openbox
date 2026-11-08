using System;
using System.Collections;
using UnityEngine;
using UnityEngine.UI;


public class ComboHandle : MonoBehaviour
{
    public Button btnNo1;
    public Button btnNo2;
    public Button btnNo3;
    public Button btnNo4;

    private Action<int> handleClickCB = null;
    public void Init(Action<int> m_handleClickCB)
    {
        handleClickCB = m_handleClickCB;
    }

    private void OnButtonClick(int comboNo)
    {
        handleClickCB?.Invoke(comboNo);
    }

    public void OnEnable()
    {
        btnNo1?.onClick.AddListener(() => { OnButtonClick(1); });
        btnNo2?.onClick.AddListener(() => { OnButtonClick(2); });
        btnNo3?.onClick.AddListener(() => { OnButtonClick(3); });
        btnNo4?.onClick.AddListener(() => { OnButtonClick(4); });
    }

    public void OnDisable()
    {
        btnNo1?.onClick.RemoveAllListeners();
        btnNo2?.onClick.RemoveAllListeners();
        btnNo3?.onClick.RemoveAllListeners();
        btnNo4?.onClick.RemoveAllListeners();
    }
}
