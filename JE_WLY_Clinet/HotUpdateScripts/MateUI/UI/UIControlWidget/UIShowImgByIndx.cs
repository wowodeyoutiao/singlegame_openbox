using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class UIShowImgByIndx : MonoBehaviour
{
    private Image curImage = null;
    public List<Sprite> SpriteCollect = new List<Sprite>();
    

    public void Init(bool m_isOn,Action<bool> m_onClickCB)
    {
        curImage = GetComponent<Image>();
        if (curImage == null)
            Debug.LogError($"{this.gameObject.name}缺少Image组件，请检查");

        if(SpriteCollect.Count == 0)
            Debug.LogError($"在SpriteCollect未发现Sprite信息");
    }


    public void ShowImage(int index)
    {
        if ((index < 0 && index >= SpriteCollect.Count))
        {
            Debug.LogError($"参数index={index} 的取值范围超出了SpriteCollect = {SpriteCollect.Count}的最大或最小值");
            return;
        }
        curImage.sprite = SpriteCollect[index];
    }

    public void OnEnable()
    {
        
    }

    public void OnDisable()
    {
        
    }
}