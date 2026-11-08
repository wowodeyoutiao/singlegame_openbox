using DG.Tweening;
using JEngine.Helper;
using UnityEngine;

public class UIHallSpecialSKillName : MonoBehaviour
{
    public void Show()
    {
        var pos = transform.localPosition;
        pos.x = -1000f;
        transform.localPosition = pos;
        gameObject.SetActive(true);

        var seq = DOTween.Sequence();
        seq.Append(transform.DOLocalMoveX(-170, 0.25f));
        seq.AppendInterval(1f);
        seq.Append(transform.DOLocalMoveX(1500, 0.5f));
        seq.AppendCallback(() =>
        {
            seq.Kill();
            seq = null;
            gameObject.SetActive(false);
            PoolHelper.Instance.RecycObj(gameObject);
        });
        seq.Play();
    }
}