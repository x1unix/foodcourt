/**
 * Class-basis for extension for components with network communication
 *
 * Class contains load status tracking and management that simplifies
 * UI state management (preloaders, errors) for components that depends
 * on some data from HTTP requests.
 *
 * <b>Note</b>: this class used <b>only</b> as component boilerplate and must
 * be extended by the component class.
 *
 * If you need to have a status for multiple resources - use <b>ResourceStatus</b> class please.
 *
 * @example 'MyComponent extends LoadStatusComponent'
 * @export
 * @class LoadStatusComponent
 */
import {ResourceStatus} from './resource-status';

export abstract class LoadStatusComponent extends ResourceStatus {}
