package connector

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"go.mau.fi/mautrix-meta/messagix"
	"go.mau.fi/mautrix-meta/messagix/cookies"
	"go.mau.fi/mautrix-meta/messagix/types"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
)

type MetaClient struct {
	Main   *MetaConnector
	client *messagix.Client
	log    zerolog.Logger
}

func cookiesFromMetadata(metadata map[string]interface{}) *cookies.Cookies {
	platform := types.Platform(metadata["platform"].(float64))

	m := make(map[string]string)
	for k, v := range metadata["cookies"].(map[string]interface{}) {
		m[k] = v.(string)
	}

	c := &cookies.Cookies{
		Platform: platform,
	}
	c.UpdateValues(m)
	return c
}

func NewMetaClient(ctx context.Context, main *MetaConnector, login *bridgev2.UserLogin) (*MetaClient, error) {
	cookies := cookiesFromMetadata(login.Metadata.Extra)

	log := login.User.Log.With().Str("component", "messagix").Logger()
	client := messagix.NewClient(cookies, log)

	return &MetaClient{
		Main:   main,
		client: client,
		log:    login.User.Log,
	}, nil
}

func (m *MetaClient) eventHandler(rawEvt any) {
	switch evt := rawEvt.(type) {
	/*
		case *messagix.Event_PublishResponse:
			user.log.Trace().Any("table", &evt.Table).Msg("Got new event")
			select {
			case user.incomingTables <- evt.Table:
			default:
				user.log.Warn().Msg("Incoming tables channel full, event order not guaranteed")
				go func() {
					user.incomingTables <- evt.Table
				}()
			}
		case *messagix.Event_Ready:
			user.log.Debug().Msg("Initial connect to Meta socket completed")
			user.metaState = status.BridgeState{StateEvent: status.StateConnected}
			user.BridgeState.Send(user.metaState)
			if initTable := user.initialTable.Swap(nil); initTable != nil {
				user.log.Debug().Msg("Sending cached initial table to handler")
				user.incomingTables <- initTable
			}
			if user.bridge.Config.Meta.Mode.IsMessenger() || user.bridge.Config.Meta.IGE2EE {
				go func() {
					err := user.connectE2EE()
					if err != nil {
						user.log.Err(err).Msg("Error connecting to e2ee")
					}
				}()
			}
			go user.BackfillLoop()
		case *messagix.Event_SocketError:
			user.log.Debug().Err(evt.Err).Msg("Disconnected from Meta socket")
			user.metaState = status.BridgeState{
				StateEvent: status.StateTransientDisconnect,
				Error:      MetaTransientDisconnect,
			}
			if evt.ConnectionAttempts > setDisconnectStateAfterConnectAttempts {
				user.BridgeState.Send(user.metaState)
			}
		case *messagix.Event_Reconnected:
			user.log.Debug().Msg("Reconnected to Meta socket")
			user.metaState = status.BridgeState{StateEvent: status.StateConnected}
			user.BridgeState.Send(user.metaState)
		case *messagix.Event_PermanentError:
			if errors.Is(evt.Err, messagix.CONNECTION_REFUSED_UNAUTHORIZED) {
				user.metaState = status.BridgeState{
					StateEvent: status.StateBadCredentials,
					Error:      MetaConnectionUnauthorized,
				}
			} else if errors.Is(evt.Err, messagix.CONNECTION_REFUSED_SERVER_UNAVAILABLE) {
				if user.bridge.Config.Meta.Mode.IsMessenger() {
					user.metaState = status.BridgeState{
						StateEvent: status.StateUnknownError,
						Error:      MetaServerUnavailable,
					}
					if user.canReconnect() {
						user.log.Debug().Msg("Doing full reconnect after server unavailable error")
						go user.FullReconnect()
					}
				} else {
					user.metaState = status.BridgeState{
						StateEvent: status.StateBadCredentials,
						Error:      IGChallengeRequiredMaybe,
					}
				}
			} else {
				user.metaState = status.BridgeState{
					StateEvent: status.StateUnknownError,
					Error:      MetaPermanentError,
					Message:    evt.Err.Error(),
				}
			}
			user.BridgeState.Send(user.metaState)
			go user.sendMarkdownBridgeAlert(context.TODO(), "Error in %s connection: %v", user.bridge.ProtocolName, evt.Err)
			user.StopBackfillLoop()
			if user.forceRefreshTimer != nil {
				user.forceRefreshTimer.Stop()
			}
	*/
	default:
		m.log.Warn().Type("event_type", evt).Msg("Unrecognized event type from messagix")
	}
}

func (m *MetaClient) Connect(ctx context.Context) error {
	// We have to call this before calling `Connect`, even if we don't use the result
	_, _, err := m.client.LoadMessagesPage()
	if err != nil {
		return fmt.Errorf("failed to load messages page: %w", err)
	}
	m.client.SetEventHandler(m.eventHandler)
	err = m.client.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to messagix: %w", err)
	}
	return nil
}

// Disconnect implements bridgev2.NetworkAPI.
func (m *MetaClient) Disconnect() {
	panic("unimplemented")
}

// GetCapabilities implements bridgev2.NetworkAPI.
func (m *MetaClient) GetCapabilities(ctx context.Context, portal *bridgev2.Portal) *bridgev2.NetworkRoomCapabilities {
	panic("unimplemented")
}

// GetChatInfo implements bridgev2.NetworkAPI.
func (m *MetaClient) GetChatInfo(ctx context.Context, portal *bridgev2.Portal) (*bridgev2.PortalInfo, error) {
	panic("unimplemented")
}

// GetUserInfo implements bridgev2.NetworkAPI.
func (m *MetaClient) GetUserInfo(ctx context.Context, ghost *bridgev2.Ghost) (*bridgev2.UserInfo, error) {
	panic("unimplemented")
}

// HandleMatrixMessage implements bridgev2.NetworkAPI.
func (m *MetaClient) HandleMatrixMessage(ctx context.Context, msg *bridgev2.MatrixMessage) (message *bridgev2.MatrixMessageResponse, err error) {
	panic("unimplemented")
}

// IsLoggedIn implements bridgev2.NetworkAPI.
func (m *MetaClient) IsLoggedIn() bool {
	panic("unimplemented")
}

// IsThisUser implements bridgev2.NetworkAPI.
func (m *MetaClient) IsThisUser(ctx context.Context, userID networkid.UserID) bool {
	panic("unimplemented")
}

// LogoutRemote implements bridgev2.NetworkAPI.
func (m *MetaClient) LogoutRemote(ctx context.Context) {
	panic("unimplemented")
}

var (
	_ bridgev2.NetworkAPI = (*MetaClient)(nil)
	// _ bridgev2.EditHandlingNetworkAPI        = (*MetaClient)(nil)
	// _ bridgev2.ReactionHandlingNetworkAPI    = (*MetaClient)(nil)
	// _ bridgev2.RedactionHandlingNetworkAPI   = (*MetaClient)(nil)
	// _ bridgev2.ReadReceiptHandlingNetworkAPI = (*MetaClient)(nil)
	// _ bridgev2.ReadReceiptHandlingNetworkAPI = (*MetaClient)(nil)
	// _ bridgev2.TypingHandlingNetworkAPI      = (*MetaClient)(nil)
	// _ bridgev2.IdentifierResolvingNetworkAPI = (*MetaClient)(nil)
	// _ bridgev2.GroupCreatingNetworkAPI       = (*MetaClient)(nil)
	// _ bridgev2.ContactListingNetworkAPI      = (*MetaClient)(nil)
)
