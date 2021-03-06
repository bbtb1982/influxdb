// Libraries
import React, {PureComponent} from 'react'
import {WithRouterProps, withRouter} from 'react-router'

import _ from 'lodash'

// Components
import {
  Form,
  Button,
  ComponentSize,
  Panel,
  IconFont,
  ComponentSpacer,
  AlignItems,
  FlexDirection,
  Gradients,
} from '@influxdata/clockface'
import {ErrorHandling} from 'src/shared/decorators/errors'

// Types
import {ButtonType} from 'src/clockface'

type Props = WithRouterProps

@ErrorHandling
class OrgProfileTab extends PureComponent<Props> {
  public render() {
    return (
      <Panel size={ComponentSize.Medium}>
        <Panel.Header title="Organization Profile" />
        <Panel.Body>
          <Form onSubmit={this.handleShowEditOverlay}>
            <Panel
              gradient={Gradients.DocScott}
              size={ComponentSize.ExtraSmall}
            >
              <Panel.Header title="Danger Zone!" />
              <Panel.Body>
                <ComponentSpacer
                  alignItems={AlignItems.Center}
                  direction={FlexDirection.Row}
                  margin={ComponentSize.Large}
                >
                  <div>
                    <h4>Rename Organization</h4>
                    <p>
                      This action can have wide-reaching unintended
                      consequences.
                    </p>
                  </div>
                  <Button
                    text="Rename"
                    icon={IconFont.Pencil}
                    type={ButtonType.Submit}
                  />
                </ComponentSpacer>
              </Panel.Body>
            </Panel>
          </Form>
        </Panel.Body>
      </Panel>
    )
  }

  private handleShowEditOverlay = () => {
    const {
      params: {orgID},
      router,
    } = this.props

    router.push(`/orgs/${orgID}/profile/rename`)
  }
}

export default withRouter<{}>(OrgProfileTab)
