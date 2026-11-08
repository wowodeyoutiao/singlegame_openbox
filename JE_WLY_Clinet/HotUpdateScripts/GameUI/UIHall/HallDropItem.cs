using System;
using UnityEngine;

public class HallDropItem : MonoBehaviour
{
    private bool moveEnable = true;

    private bool waitDie = false;
    private float dieTimeInterval;
    private Action callBack;

    private Vector3 startPos = Vector3.zero;
    private Vector3 endPos = Vector3.zero;
    private float statTime;

    private float timeInterval = 0f; // 每段时间

    private float[] timeScales = new float[] { 6f, 20f, 30 };
    private float timeScale;

    private float[] dampings = new float[] { 0.6f, 0.1f, 0.1f}; // 反弹衰减系数
    private float damping = 1;

    private float bounceHight = 80f;
    private float curHight = 0f;

    private int totalBounceCount = 3; // 反弹次数
    private int curBounceCount = 0;
    private float curN = 1;

    public void Reset()
    {
        transform.localPosition = Vector3.zero;
        transform.localScale = Vector3.one;

        moveEnable = false;
        damping = 1;
        curHight = 0;
        curBounceCount = 0;
        curN = 1;
    }

    public void Set(Vector3 sPos, Vector3 ePos, Action c)
    {
        startPos = sPos;
        endPos = ePos;
        statTime = Time.time;

        moveEnable = true;
        timeScale = timeScales[0];

        callBack = c;

        timeInterval = 0f;
        waitDie = false;
        gameObject.SetActive(true);
    }

    void Update()
    {
        if(waitDie)
        {
            dieTimeInterval += Time.deltaTime;

            if (dieTimeInterval > 0.8f)
            {
                waitDie = false;
                callBack?.Invoke();
            }
        }

        if (!moveEnable) return;

        var n = Mathf.Sin(timeInterval * Mathf.PI);
        curHight = bounceHight * n * damping;

        if(n * curN < 0)
        {
            curN = n;
            // 改变方向
            curBounceCount++;
            if (curBounceCount == totalBounceCount)
            {
                moveEnable = false;

                waitDie = true;
                dieTimeInterval = 0;
            }

            timeScale = timeScales[curBounceCount - 1];

            damping = dampings[curBounceCount - 1];
            if (n < 0) damping *= -1;
        }

        var pos = Vector3.Lerp(startPos, endPos, (Time.time - statTime) * 2.5f);
        pos.y += curHight;
        transform.localPosition = pos;

        timeInterval += Time.deltaTime * timeScale; //2 0.95s; 3 0.6s
    }
}