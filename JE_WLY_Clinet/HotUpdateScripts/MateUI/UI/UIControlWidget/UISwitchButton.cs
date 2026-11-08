using System;
using System.Collections;
using UnityEngine;
using UnityEngine.UI;

public class UISwitchButton : MonoBehaviour
{
    public GameObject goOn;
    public GameObject goOff;

    public Button btnClick;

    private bool isOn = false;
    private Action<bool> onClickCB = null;

    public bool IsOn { get => isOn;}

    public void Init(bool m_isOn,Action<bool> m_onClickCB)
    {
        setState(m_isOn);
        onClickCB = m_onClickCB;
    }


    private void setState(bool m_isOn)
    {
        isOn = m_isOn;

        goOn?.SetActive(isOn);
        goOff?.SetActive(!isOn);

        onClickCB?.Invoke(isOn);
    }

    private void OnButtonClick()
    {
        setState(!isOn);
    }

    public void OnEnable()
    {
        btnClick?.onClick.AddListener(OnButtonClick);
    }

    public void OnDisable()
    {
        btnClick?.onClick.RemoveAllListeners();
    }
}