using Cysharp.Threading.Tasks;
using JEngine.Helper;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class UIHallBg : MonoBehaviour
{
    private List<Image> bgImgs = new List<Image>();

    private float imgWidth = 720f;

    private bool isBossLevel;

    private int oldBgCount = 0;
    private int bgCount = 8;
    private int bgSpriteIndex; // 背景的图片资源

    private string norPreName;
    private string bossPreName;

    private bool isMoving;
    private float moveSpeed = 5f; // 移动速度
    private float speedScale; // 远景、近景移动速度不同
    private float timeScale = 1f;

    private void LoadBgImgs()
    {
        bgImgs.Clear();
        for (int i = 0; i < transform.childCount; i++)
            bgImgs.Add(transform.GetChild(i).GetComponent<Image>());
    }

    public void OnInit(string norPreName, string bossPreName, float speedScale)
    {
        this.norPreName = norPreName;
        this.bossPreName = bossPreName;
        this.speedScale = speedScale;

        isMoving = false;
    }

    public void RefreshBg(bool isBossLevel)
    {
        bgCount = isBossLevel? 2 : 8;
        this.isBossLevel = isBossLevel;

        if (oldBgCount == bgCount)
            return;

        if (bgImgs.Count == 0)
            LoadBgImgs();

        oldBgCount = bgCount;

        var index = 0;
        for (int i = 0; i < bgImgs.Count; i++)
        {
            var bgImg = bgImgs[i];
            bgImg.gameObject.SetActive(true);

            bgImg.transform.localPosition = new Vector3(imgWidth * (i - 1), 0, 0);

            var rect = bgImg.transform.GetComponent<RectTransform>();
            rect.sizeDelta = new Vector2(imgWidth, rect.sizeDelta.y);

            index = i;
            if (i >= bgCount)
                index -= bgCount;

            if (isBossLevel)
                ResouseHelper.Instance.LoadSprite($"Common/Texture/{bossPreName}{index + 1}", bgImg).Forget();
            else
                ResouseHelper.Instance.LoadSprite($"Common/Texture/{norPreName}{index + 1}", bgImg).Forget();

            bgSpriteIndex = index;
        }
    }

    public void StartBgMove(float timeScale = 1)
    {
        isMoving = true;
        this.timeScale = timeScale;
    }

    public void StopBgMove()
    {
        isMoving = false;
    }

    private void FixedUpdate()
    {
        if (!isMoving)
            return;

        for (int i = 0; i < bgImgs.Count; i++)
        {
            var pos = bgImgs[i].transform.localPosition;
            pos.x -= moveSpeed * timeScale * speedScale;
            if (pos.x <= -imgWidth - 100)
            {
                pos.x += imgWidth * 3;
                ChangeBgSprite(bgImgs[i]);
            }
            bgImgs[i].GetComponent<RectTransform>().anchoredPosition = pos;
        }
    }

    private void ChangeBgSprite(Image img)
    {
        bgSpriteIndex++;
        if (bgSpriteIndex >= bgCount)
            bgSpriteIndex = 0;

        if (isBossLevel)
            ResouseHelper.Instance.LoadSprite($"Common/Texture/{bossPreName}{bgSpriteIndex + 1}", img).Forget();
        else
            ResouseHelper.Instance.LoadSprite($"Common/Texture/{norPreName}{bgSpriteIndex + 1}", img).Forget();
    }
}