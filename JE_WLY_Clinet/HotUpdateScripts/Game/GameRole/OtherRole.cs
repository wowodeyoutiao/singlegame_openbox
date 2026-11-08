using HotUpdateScripts.Game;
using HotUpdateScripts.ProtoM;
using pb;
using UnityEngine;

public class OtherRole : Role
{
    public Direction _joyDir;

    private long _targetID;//目标攻击玩家ID
    public long targetID
    {
        get => _targetID;

        set
        {
            if (_targetID == value)
                return;
            _targetID = value;

        }
    }


    public override void OnInit()
    {
        base.OnInit();
        _joyDir = Direction.NONE;

        PlayerInfo _playerInfo = GameRoomM.Instance.getPlayerInfoByID(uID);
        Direction direction = default;
        if (_playerInfo != null)
        {
            pointNow = _playerInfo.Location;
            direction = _playerInfo.Direction;
            wrap(new Vector3((float)pointNow.X, (float)pointNow.Y, SceneManager.GetOffZ(pointNow.X, pointNow.Y)));
        }
        else
        {
            ObjectInfo objectInfo = GameRoomM.Instance.GetOneObjInfoByID(uID);
            if (objectInfo != null)
            {
                direction = objectInfo.Direction;
                pointNow = objectInfo.Location;
            }
        }

        if (direction != default)
        {
            wrap(new Vector3((float)pointNow.X, (float)pointNow.Y, SceneManager.GetOffZ(pointNow.X, pointNow.Y)));
            setRotation(direction);
        }

    }

}


