using HotUpdateScripts.Game;
using UnityEngine;

public class FollowTarget : MonoBehaviour
{
    //public Transform hero;
    public static FollowTarget _FollowTarget;
    private Vector3 offset;

    public float damping = 5;

    private void Awake()
    {
        _FollowTarget = this;
    }

    void Start()
    {
        //offset = new Vector3(0,15.9f, -13);
        offset = new Vector3(0, 4f, -17);
        Debug.Log("followtarget offset:" + offset);
    }

    private void LateUpdate‌()
    {
        if (offset == Vector3.zero)
            return;

        if (MainRole._MainRole == null)
        {
            Debug.Log("MainRole._MainRole == null return");
            return;
        }

        //if (MainRole._MainRole.isJoyWalking)
        //    return;

        if (damping > 0)
        {

            Vector3 targetPos = MainRole._MainRole.transform.position;
            if (targetPos.x <= SceneManager.MinX)
                targetPos.x = SceneManager.MinX;
            else if (targetPos.x >= SceneManager.MaxX)
                targetPos.x = SceneManager.MaxX;

            if (targetPos.y <= SceneManager.MinY)
                targetPos.y = SceneManager.MinY;
            else if (targetPos.y >= SceneManager.MaxY)
                targetPos.y = SceneManager.MaxY;

            Vector3 desiredPosition = targetPos + offset;

            Vector3 position = Vector3.Lerp(transform.position, desiredPosition, Time.deltaTime * damping);

            transform.position = position;

            //  transform.LookAt(hero.position);

        }
        else
        {

            Vector3 desiredPosition = MainRole._MainRole.transform.position + offset;
            transform.position = desiredPosition;

        }
    }

    public void SetPos(Vector3 basePos)
    {
        transform.position = basePos + offset;
    }
}
